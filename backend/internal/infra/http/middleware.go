package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-shop-app-backend/internal/infra/auth"
)

const (
	ctxUserIDKey   = "userID"
	ctxUserRoleKey = "userRole"
)

func AuthMiddleware(jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_authorization_header",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_authorization_header",
				"message": "Authorization header must be in format: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "empty_token",
				"message": "Bearer token is empty",
			})
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxUserRoleKey, claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, ok := c.Get(ctxUserRoleKey)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "user role not found in context",
			})
			c.Abort()
			return
		}

		role, _ := roleVal.(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (int64, bool) {
	val, ok := c.Get(ctxUserIDKey)
	if !ok {
		return 0, false
	}

	id, ok := val.(int64)
	if !ok {
		return 0, false
	}

	return id, true
}

func GetUserRole(c *gin.Context) (string, bool) {
	val, ok := c.Get(ctxUserRoleKey)
	if !ok {
		return "", false
	}

	role, ok := val.(string)
	if !ok {
		return "", false
	}

	return role, true
}
