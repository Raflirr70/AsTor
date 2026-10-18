package middleware

import (
	"astor/backend/pkg/jwt"
	"astor/backend/pkg/response"
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func AuthMiddleware(jwtSecret string, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := jwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		ctx := context.Background()
		blacklistKey := fmt.Sprintf("blacklist:%s", tokenString)
		val, err := redisClient.Get(ctx, blacklistKey).Result()
		if err == nil && val == "1" {
			response.Unauthorized(c, "token has been revoked (logged out)")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", tokenString)
		c.Next()
	}
}
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			response.Forbidden(c, "invalid role")
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "access denied")
		c.Abort()
	}
}
