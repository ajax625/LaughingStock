package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/config"
	"github.com/laughingstock/v1/internal/api"
	"github.com/laughingstock/v1/internal/broker"
	"github.com/laughingstock/v1/internal/notify"
	"github.com/laughingstock/v1/internal/store"
	"github.com/laughingstock/v1/internal/tradex"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("db connect failed", zap.Error(err))
	}
	defer db.Close()

	hub := notify.NewHub()

	tg, err := notify.NewTelegramBot(cfg.TelegramBotToken, logger)
	if err != nil {
		logger.Fatal("telegram init failed", zap.Error(err))
	}

	tx := tradex.New(cfg.TradeXBaseURL, cfg.TradeXAPIKey)

	handler := api.NewHandler(cfg, db, hub, tg, tx, logger)

	router := gin.New()
	router.Use(gin.Recovery())
	handler.Register(router)

	// Start Telegram callback listener — handles inline Execute button presses
	if tg != nil {
		go tg.StartCallbackListener(func(chatID int64, symbol, direction string) {
			ctx := context.Background()

			user, err := db.GetUserByTelegramChatID(ctx, strconv.FormatInt(chatID, 10))
			if err != nil || user == nil {
				logger.Warn("telegram callback: user not found", zap.Int64("chat_id", chatID))
				tg.SendMessage(chatID, "Could not find your account. Please check your Telegram Chat ID in settings.")
				return
			}
			if user.AlpacaKey == "" {
				tg.SendMessage(chatID, "Alpaca credentials not configured. Please add them in your LaughingStock account.")
				return
			}

			side := "buy"
			if direction == "SHORT" {
				side = "sell"
			}

			result, err := broker.PlaceMarketOrder(ctx, user.AlpacaKey, user.AlpacaSecret, symbol, side, 1)
			if err != nil {
				logger.Error("telegram callback: order failed",
					zap.String("user_id", user.ID),
					zap.String("symbol", symbol),
					zap.Error(err),
				)
				tg.SendMessage(chatID, fmt.Sprintf("Order failed for %s: %s", symbol, err.Error()))
				return
			}

			logger.Info("telegram callback: order placed",
				zap.String("user_id", user.ID),
				zap.String("order_id", result.ID),
				zap.String("symbol", symbol),
				zap.String("side", side),
			)
			tg.SendMessage(chatID, fmt.Sprintf("✅ Order placed: %s %s x1\nOrder ID: %s\nStatus: %s",
				side, symbol, result.ID, result.Status,
			))
		})
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("laughingstock starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
	logger.Info("laughingstock stopped")
}
