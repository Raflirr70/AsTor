package middleware

import (
	"astor/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// TODO: log r
				response.Internal(c, "Terjadi kesalahan internal")
				c.Abort()
			}
		}()
		c.Next()
		// kalo ada error yg diset via c.Error() bisa ditangani sini
	}
}
