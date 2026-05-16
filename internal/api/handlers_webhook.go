package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Webhook receives a TradeX signal payload and routes it to the user
// via WebSocket (if online) or Telegram (if offline).
func (h *Handler) Webhook(c *gin.Context) {
	token := c.Param("token")

	user, err := h.store.GetUserByToken(c.Request.Context(), token)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown token"})
		return
	}

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil || len(payload) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty payload"})
		return
	}

	if !json.Valid(payload) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if h.hub.IsOnline(user.ID) {
		h.hub.Push(user.ID, payload)
		h.log.Info("webhook: pushed via WebSocket", zap.String("user_id", user.ID))
	} else if h.telegram != nil && user.TelegramChatID != "" {
		chatID, _ := strconv.ParseInt(user.TelegramChatID, 10, 64)
		if err := h.telegram.SendSignalAlert(chatID, payload); err != nil {
			h.log.Error("webhook: telegram send failed", zap.String("user_id", user.ID), zap.Error(err))
		} else {
			h.log.Info("webhook: sent via Telegram", zap.String("user_id", user.ID))
		}
	} else {
		h.log.Warn("webhook: user offline and no Telegram configured", zap.String("user_id", user.ID))
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ServeWS upgrades the connection to WebSocket for real-time signal delivery.
func (h *Handler) ServeWS(c *gin.Context) {
	h.hub.ServeWS(c.Writer, c.Request, currentUserID(c))
}
