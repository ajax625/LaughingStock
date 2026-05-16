package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/broker"
)

func (h *Handler) GetQuotes(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.store.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load user"})
		return
	}
	if user.AlpacaKey == "" || user.AlpacaSecret == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	rows, err := h.store.GetUserSymbols(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load symbols"})
		return
	}

	syms := make([]string, 0, len(rows))
	for _, s := range rows {
		if s.Active {
			syms = append(syms, s.Symbol)
		}
	}

	prices, err := broker.GetLatestPrices(user.AlpacaKey, user.AlpacaSecret, syms)
	if err != nil {
		h.log.Sugar().Warnf("quotes: alpaca fetch failed: %v", err)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, prices)
}
