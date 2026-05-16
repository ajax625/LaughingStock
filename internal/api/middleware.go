package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/laughingstock/v1/internal/auth"
)

const userIDKey = "userID"
const userEmailKey = "userEmail"

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := tokenFromRequest(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		claims, err := auth.ValidateToken(token, h.cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.Email)
		c.Next()
	}
}

func tokenFromRequest(c *gin.Context) string {
	// 1. HttpOnly cookie
	if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
		return cookie
	}
	// 2. Authorization: Bearer <token>
	if h := c.GetHeader("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}

func currentUserID(c *gin.Context) string {
	id, _ := c.Get(userIDKey)
	s, _ := id.(string)
	return s
}
