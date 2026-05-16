package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/auth"
	"github.com/laughingstock/v1/internal/store"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, _ := h.store.GetUserByEmail(c.Request.Context(), body.Email)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	id := fmt.Sprintf("usr-%d", time.Now().UnixNano())
	token := generateToken()

	u := &store.User{
		ID:           id,
		Email:        body.Email,
		PasswordHash: string(hash),
		WebhookToken: token,
	}
	if err := h.store.CreateUser(c.Request.Context(), u); err != nil {
		h.log.Error("register: create user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "email": body.Email, "webhook_token": token})
}

func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.store.GetUserByEmail(c.Request.Context(), body.Email)
	if err != nil || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	jwtToken, err := auth.IssueToken(u.ID, u.Email, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.SetCookie("token", jwtToken, int(24*time.Hour/time.Second), "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) GetMe(c *gin.Context) {
	u, err := h.store.GetUserByID(c.Request.Context(), currentUserID(c))
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":               u.ID,
		"email":            u.Email,
		"telegram_chat_id": u.TelegramChatID,
		"has_alpaca":       u.AlpacaKey != "",
		"webhook_token":    u.WebhookToken,
	})
}

func (h *Handler) UpdateMe(c *gin.Context) {
	var body struct {
		TelegramChatID string `json:"telegram_chat_id"`
		AlpacaKey      string `json:"alpaca_key"`
		AlpacaSecret   string `json:"alpaca_secret"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.UpdateUserProfile(c.Request.Context(), currentUserID(c), body.TelegramChatID, body.AlpacaKey, body.AlpacaSecret); err != nil {
		h.log.Error("update me failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
