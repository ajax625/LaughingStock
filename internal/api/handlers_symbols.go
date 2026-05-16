package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/store"
	"github.com/laughingstock/v1/internal/tradex"
	"go.uber.org/zap"
)

func (h *Handler) GetSymbols(c *gin.Context) {
	syms, err := h.store.GetUserSymbols(c.Request.Context(), currentUserID(c))
	if err != nil {
		h.log.Error("get symbols failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve symbols"})
		return
	}
	if syms == nil {
		syms = []*store.UserSymbol{}
	}
	c.JSON(http.StatusOK, gin.H{"symbols": syms})
}

func (h *Handler) AddSymbol(c *gin.Context) {
	var body struct {
		Symbol        string  `json:"symbol" binding:"required"`
		TraderType    string  `json:"trader_type"`
		MinROI        float64 `json:"min_roi"`
		MinRR         float64 `json:"min_rr"`
		MinConfidence float64 `json:"min_confidence"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.TraderType == "" {
		body.TraderType = "day"
	}
	if body.MinROI == 0 {
		body.MinROI = 0.5
	}
	if body.MinRR == 0 {
		body.MinRR = 2.0
	}
	if body.MinConfidence == 0 {
		body.MinConfidence = 0.6
	}

	userID := currentUserID(c)
	user, err := h.store.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	webhookURL := fmt.Sprintf("%s/webhook/%s", h.cfg.LaughingStockBaseURL, user.WebhookToken)

	tradexSubID, err := h.tradex.CreateSubscription(c.Request.Context(), tradex.CreateSubscriptionRequest{
		Name:          fmt.Sprintf("%s-%s", user.Email, body.Symbol),
		WebhookURL:    webhookURL,
		Symbol:        body.Symbol,
		TraderType:    body.TraderType,
		MinROI:        body.MinROI,
		MinRR:         body.MinRR,
		MinConfidence: body.MinConfidence,
	})
	if err != nil {
		h.log.Error("add symbol: tradex register failed", zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to register with TradeX"})
		return
	}

	id := fmt.Sprintf("sym-%d", time.Now().UnixNano())
	us := &store.UserSymbol{
		ID:                   id,
		UserID:               userID,
		Symbol:               body.Symbol,
		TraderType:           body.TraderType,
		MinROI:               body.MinROI,
		MinRR:                body.MinRR,
		MinConfidence:        body.MinConfidence,
		TradeXSubscriptionID: tradexSubID,
	}
	if err := h.store.CreateUserSymbol(c.Request.Context(), us); err != nil {
		h.log.Error("add symbol: store failed", zap.Error(err))
		_ = h.tradex.DeleteSubscription(c.Request.Context(), tradexSubID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save symbol"})
		return
	}

	// Trigger fundamentals research async — don't block the response.
	sym := body.Symbol
	go func() {
		if _, err := h.tradex.TriggerDeepResearch(context.Background(), sym); err != nil {
			h.log.Warn("add symbol: deep research trigger failed", zap.String("symbol", sym), zap.Error(err))
		}
	}()

	c.JSON(http.StatusCreated, gin.H{"id": id, "symbol": body.Symbol, "tradex_subscription_id": tradexSubID})
}

func (h *Handler) UpdateSymbol(c *gin.Context) {
	id := c.Param("id")
	userID := currentUserID(c)

	var body struct {
		MinROI        float64 `json:"min_roi"`
		MinRR         float64 `json:"min_rr"`
		MinConfidence float64 `json:"min_confidence"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	us, err := h.store.GetUserSymbolByID(c.Request.Context(), id, userID)
	if err != nil || us == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "symbol not found"})
		return
	}

	if us.TradeXSubscriptionID != "" {
		if err := h.tradex.UpdateSubscription(c.Request.Context(), us.TradeXSubscriptionID, tradex.UpdateSubscriptionRequest{
			MinROI:        body.MinROI,
			MinRR:         body.MinRR,
			MinConfidence: body.MinConfidence,
		}); err != nil {
			h.log.Warn("update symbol: tradex update failed", zap.Error(err))
		}
	}

	if err := h.store.UpdateUserSymbol(c.Request.Context(), id, userID, body.MinROI, body.MinRR, body.MinConfidence); err != nil {
		h.log.Error("update symbol: store failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update symbol"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated": id})
}

func (h *Handler) DeleteSymbol(c *gin.Context) {
	id := c.Param("id")
	userID := currentUserID(c)

	us, err := h.store.GetUserSymbolByID(c.Request.Context(), id, userID)
	if err != nil || us == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "symbol not found"})
		return
	}

	if us.TradeXSubscriptionID != "" {
		if err := h.tradex.DeleteSubscription(c.Request.Context(), us.TradeXSubscriptionID); err != nil {
			h.log.Warn("delete symbol: tradex deregister failed", zap.Error(err))
		}
	}

	if err := h.store.DeactivateUserSymbol(c.Request.Context(), id, userID); err != nil {
		h.log.Error("delete symbol: store failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete symbol"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
