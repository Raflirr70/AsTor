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
	loanHandler *handler.LoanHandler,
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
		admin.PUT("/labkom/:id", laboratoryHandler.Update)
		admin.GET("/labkom/:id", laboratoryHandler.GetByID)
		admin.DELETE("/labkom/:id", laboratoryHandler.Delete)
		admin.PATCH("/labkom/:id/restore", laboratoryHandler.Restore)

		admin.GET("/specifications", specificationHandler.GetAll)
		admin.POST("/specifications", specificationHandler.Create)
		admin.GET("/specification/:id", specificationHandler.GetByID)
		admin.PUT("/specification/:id", specificationHandler.Update)
		admin.DELETE("/specification/:id", specificationHandler.Delete)

		admin.GET("/categories", categoryHandler.GetAll)
		admin.POST("/categories", categoryHandler.Create)
		admin.GET("/category/:id", categoryHandler.GetByID)
		admin.PUT("/category/:id", categoryHandler.Update)
		admin.DELETE("/category/:id", categoryHandler.Delete)

		admin.GET("/assets", assetHandler.GetAll)
		admin.POST("/assets", assetHandler.Create)
		admin.GET("/asset/:id", assetHandler.GetByID)
		admin.PUT("/asset/:id", assetHandler.Update)
		admin.DELETE("/asset/:id", assetHandler.Delete)
		admin.PATCH("/asset/:id/checking", assetHandler.Checking)

		admin.GET("/computers", computerHandler.GetAll)
		admin.POST("/computers", computerHandler.Create)
		admin.GET("/computer/:id", computerHandler.GetByID)
		admin.PUT("/computer/:id", computerHandler.Update)
		admin.DELETE("/computer/:id", computerHandler.Delete)
	}
	gudang := r.Group("/gudang")
	gudang.Use(middleware.AuthMiddleware(jwt, redisClient))
	gudang.Use(middleware.RoleMiddleware("PJ Gudang"))
	{
		gudang.PATCH("/loan/:id/approved", loanHandler.Approved)
	}

	protected := r.Group("/asisten")
	protected.Use(middleware.AuthMiddleware(jwt, redisClient))
	{
		protected.GET("/loans", loanHandler.GetAll)
		protected.POST("/loans", loanHandler.Create)
		protected.PATCH("/loan/:id/returned", loanHandler.Returned)
	}

	return r
}
