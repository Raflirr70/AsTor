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
	specificationHandler *handler.SpecificationHandler,
	categoryHandler *handler.CategoryHandler,
	assetHandler *handler.AssetHandler,
	computerHandler *handler.ComputerHandler,
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

		admin.GET("/specificaions", specificationHandler.GetAll)
		admin.POST("/specificaions", specificationHandler.Create)
		admin.PUT("/specificaion/:id", specificationHandler.Update)
		admin.DELETE("/specificaion/:id", specificationHandler.Delete)

		admin.GET("/categories", categoryHandler.GetAll)
		admin.POST("/categories", categoryHandler.Create)
		admin.PUT("/categorie/:id", categoryHandler.Update)
		admin.DELETE("/categorie/:id", categoryHandler.Delete)

		admin.GET("/assets", assetHandler.GetAll)
		admin.POST("/assets", assetHandler.Create)
		admin.PUT("/asset/:id", assetHandler.Update)
		admin.DELETE("/asset/:id", assetHandler.Delete)
		admin.DELETE("/asset/:id/checking", assetHandler.Checking)

		admin.GET("/Computers", computerHandler.GetAll)
		admin.POST("/Computers", computerHandler.Create)
		admin.PUT("/Computer/:id", computerHandler.Update)
	}

	return r
}
