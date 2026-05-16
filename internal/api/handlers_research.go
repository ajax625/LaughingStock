package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetSymbolResearch(c *gin.Context) {
	id := c.Param("id")
	userID := currentUserID(c)

	us, err := h.store.GetUserSymbolByID(c.Request.Context(), id, userID)
	if err != nil || us == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "symbol not found"})
		return
	}

	job, err := h.tradex.GetLatestResearch(c.Request.Context(), us.Symbol)
	if err != nil {
		h.log.Warn("get research: tradex fetch failed", zap.String("symbol", us.Symbol), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"status": "unavailable"})
		return
	}
	if job == nil {
		c.JSON(http.StatusOK, gin.H{"status": "none"})
		return
	}

	c.JSON(http.StatusOK, job)
}
