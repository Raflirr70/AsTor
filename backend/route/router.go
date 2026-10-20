package route

import (
	"astor/backend/handler"
	"astor/backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupRoutes(
	authHandler *handler.AuthHandler,
	laboratoryHandler *handler.LaboratoryHandler,
	jwt string,
	redisClient *redis.Client,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"massage": "API is Running"})
	})

	r.POST("/login", authHandler.Login)
	r.POST("/logout", middleware.AuthMiddleware(jwt, redisClient), authHandler.Logout)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwt, redisClient))
	admin.Use(middleware.RoleMiddleware("Admin"))
	{
		admin.POST("/users", authHandler.Register)

		admin.POST("/labkoms", laboratoryHandler.Create)
		admin.GET("/labkoms", laboratoryHandler.GetAll)
		admin.PUT("/labkoms/:id", laboratoryHandler.Update)
		admin.DELETE("/labkom/:id", laboratoryHandler.Delete)
		admin.PATCH("/labkoms/:id/restore", laboratoryHandler.Restore)
	}

	return r
}
