package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/broker"
	"go.uber.org/zap"
)

func (h *Handler) ExecuteTrade(c *gin.Context) {
	var body struct {
		Symbol string  `json:"symbol" binding:"required"`
		Side   string  `json:"side" binding:"required"` // "buy" or "sell"
		Qty    float64 `json:"qty" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Side != "buy" && body.Side != "sell" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "side must be 'buy' or 'sell'"})
		return
	}

	user, err := h.store.GetUserByID(c.Request.Context(), currentUserID(c))
	if err != nil || user == nil || user.AlpacaKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alpaca credentials not configured"})
		return
	}

	result, err := broker.PlaceMarketOrder(c.Request.Context(), user.AlpacaKey, user.AlpacaSecret, body.Symbol, body.Side, body.Qty)
	if err != nil {
		h.log.Error("execute trade failed", zap.String("user_id", user.ID), zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "order placement failed"})
		return
	}

	h.log.Info("trade executed",
		zap.String("user_id", user.ID),
		zap.String("symbol", result.Symbol),
		zap.String("side", result.Side),
		zap.Float64("qty", result.Qty),
		zap.String("order_id", result.ID),
	)

	c.JSON(http.StatusOK, gin.H{
		"order_id": result.ID,
		"symbol":   result.Symbol,
		"side":     result.Side,
		"qty":      result.Qty,
		"status":   result.Status,
	})
}
