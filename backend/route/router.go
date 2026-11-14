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
	assetReportHandler *handler.AssetReportHandler,
	jwt string,
	redisClient *redis.Client,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

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
		admin.PUT("/labkom/:id", laboratoryHandler.Update)
		admin.DELETE("/labkom/:id", laboratoryHandler.Delete)
		admin.PATCH("/labkom/:id/restore", laboratoryHandler.Restore)
	}
	gudang := r.Group("/gudang")
	gudang.Use(middleware.AuthMiddleware(jwt, redisClient))
	gudang.Use(middleware.RoleMiddleware("PJ Gudang"))
	{
		gudang.GET("/specifications", specificationHandler.GetAll)
		gudang.POST("/specifications", specificationHandler.Create)
		gudang.GET("/specification/:id", specificationHandler.GetByID)
		gudang.PUT("/specification/:id", specificationHandler.Update)
		gudang.DELETE("/specification/:id", specificationHandler.Delete)

		gudang.GET("/categories", categoryHandler.GetAll)
		gudang.POST("/categories", categoryHandler.Create)
		gudang.GET("/category/:id", categoryHandler.GetByID)
		gudang.PUT("/category/:id", categoryHandler.Update)
		gudang.DELETE("/category/:id", categoryHandler.Delete)

		gudang.GET("/assets", assetHandler.GetAll)
		gudang.POST("/assets", assetHandler.Create)
		gudang.GET("/asset/:id", assetHandler.GetByID)
		gudang.PUT("/asset/:id", assetHandler.Update)
		gudang.DELETE("/asset/:id", assetHandler.Delete)
		gudang.PATCH("/asset/:id/checking", assetHandler.Checking)

		gudang.GET("/computers", computerHandler.GetAll)
		gudang.POST("/computers", computerHandler.Create)
		gudang.GET("/computer/:id", computerHandler.GetByID)
		gudang.PUT("/computer/:id", computerHandler.Update)
		gudang.DELETE("/computer/:id", computerHandler.Delete)

		gudang.PATCH("/loan/:id/approved", loanHandler.Approved)

		gudang.GET("/labkoms", laboratoryHandler.GetAll)
		gudang.GET("/labkom/:id", laboratoryHandler.GetByID)

	}

	labkom := r.Group("/lab")
	labkom.Use(middleware.AuthMiddleware(jwt, redisClient))
	labkom.Use(middleware.RoleMiddleware("PJ Labkom"))
	{

	}

	protected := r.Group("/asisten")
	protected.Use(middleware.AuthMiddleware(jwt, redisClient))
	{
		protected.GET("/loans", loanHandler.GetAll)
		protected.POST("/loans", loanHandler.Create)
		protected.PATCH("/loan/:id/returned", loanHandler.Returned)

		protected.POST("/report", assetReportHandler.Create)
		protected.GET("/report", assetReportHandler.GetAll)
		protected.GET("/report/:id", assetReportHandler.GetByID)
		protected.PATCH("/report/:id", assetReportHandler.Update)

	}

	return r
}
