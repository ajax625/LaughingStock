package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/tradex"
	"go.uber.org/zap"
)

// GetUserSignals returns the latest TradeX signal for each of the user's symbols.
func (h *Handler) GetUserSignals(c *gin.Context) {
	syms, err := h.store.GetUserSymbols(c.Request.Context(), currentUserID(c))
	if err != nil {
		h.log.Error("get user signals: fetch symbols failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch symbols"})
		return
	}

	type result struct {
		Symbol string          `json:"symbol"`
		Signal *tradex.Signal  `json:"signal"`
	}

	results := make([]result, len(syms))
	var wg sync.WaitGroup
	for i, us := range syms {
		wg.Add(1)
		go func(i int, symbol string) {
			defer wg.Done()
			sig, _ := h.tradex.GetLatestSignal(c.Request.Context(), symbol)
			results[i] = result{Symbol: symbol, Signal: sig}
		}(i, us.Symbol)
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"signals": results})
}
