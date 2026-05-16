package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/config"
	"github.com/laughingstock/v1/internal/notify"
	"github.com/laughingstock/v1/internal/store"
	"github.com/laughingstock/v1/internal/tradex"
	"go.uber.org/zap"
)

type Handler struct {
	cfg      *config.Config
	store    *store.Store
	hub      *notify.Hub
	telegram *notify.TelegramBot
	tradex   *tradex.Client
	log      *zap.Logger
}

func NewHandler(cfg *config.Config, s *store.Store, hub *notify.Hub, tg *notify.TelegramBot, tx *tradex.Client, log *zap.Logger) *Handler {
	return &Handler{cfg: cfg, store: s, hub: hub, telegram: tg, tradex: tx, log: log}
}

func (h *Handler) Register(r *gin.Engine) {
	r.Use(corsMiddleware(h.cfg.AllowedOrigins))

	// Public routes
	r.POST("/auth/register", h.RegisterUser)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)

	// Webhook receiver — public but token-gated
	r.POST("/webhook/:token", h.Webhook)

	// Serve React frontend
	r.StaticFS("/app", http.Dir("web/dist"))
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/app") })

	// Authenticated routes
	auth := r.Group("/", h.AuthMiddleware())
	{
		auth.GET("/me", h.GetMe)
		auth.PATCH("/me", h.UpdateMe)

		auth.GET("/me/symbols", h.GetSymbols)
		auth.POST("/me/symbols", h.AddSymbol)
		auth.PATCH("/me/symbols/:id", h.UpdateSymbol)
		auth.DELETE("/me/symbols/:id", h.DeleteSymbol)
		auth.GET("/me/symbols/:id/research", h.GetSymbolResearch)

		auth.GET("/me/signals", h.GetUserSignals)
		auth.GET("/me/quotes", h.GetQuotes)
		auth.POST("/me/trade", h.ExecuteTrade)

		auth.GET("/ws", h.ServeWS)
	}
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
