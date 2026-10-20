package main

import (
	"context"
	"log"

	"astor/backend/config"
	"astor/backend/handler"
	"astor/backend/middleware"
	"astor/backend/pkg/database"
	"astor/backend/pkg/logger"
	"astor/backend/pkg/ngrok"
	"astor/backend/pkg/redis"
	"astor/backend/repository"
	router "astor/backend/route"
	"astor/backend/usecase"

	"github.com/joho/godotenv"
)

func main() {
	logger.Init()
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using system env")
	}

	// 1. load configuration
	cfg := config.LoadConfig()

	// 2. connect Database
	db, err := database.ConnectDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}
	defer sqlDB.Close()
	logger.Info.Printf("PostgreSQL Connected")

	// 2. connect Database
	fwd, err := ngrok.ConnectNgrok(context.Background(), cfg.Ngrok)
	if err != nil {
		log.Fatalf("failed to connect ngrok: %v", err)
	}
	defer fwd.Close()

	redisClient := redis.NewRedisClient(cfg.Redis)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	logger.Info.Printf("Redis Connected")

	// 3. Setup Dependency
	//repo
	txRepo := repository.NewTransactionRepositoryImpl(db)
	activityLogRepo := repository.NewActivityLogRepositoryImpl(db)
	authRepo := repository.NewAuthRepositoryImpl(db)
	locationRepo := repository.NewLocationRepositoryImpl(db)
	laboratoryRepo := repository.NewLaboratoryRepositoryImpl(db)
	specificationRepo := repository.NewSpecificationRepositoryImpl(db)
	categoryRepo := repository.NewCategoryRepositoryImpl(db)
	catSpecRepo := repository.NewCategorySpecificationRepositoryImpl(db)
	//usecase
	authUC := usecase.NewAuthUsecaseImpl(authRepo, activityLogRepo, txRepo, cfg.JWT, redisClient)
	laboratoryUC := usecase.NewLaboratoryUsecaseImpl(laboratoryRepo, locationRepo, activityLogRepo, txRepo)
	specificationUC := usecase.NewSpecificationUsecaseImpl(specificationRepo, activityLogRepo, txRepo)
	categoryUC := usecase.NewCategoryUsecaseImpl(categoryRepo, catSpecRepo, activityLogRepo, txRepo)
	//handler
	authHandler := handler.NewAuthHandler(authUC)
	laboratoryHandler := handler.NewLaboratoryHandler(laboratoryUC)
	specificationHandler := handler.NewSpecificationHandler(specificationUC)
	categoryHandler := handler.NewCategoryHandler(categoryUC)

	// 3. Setup Router
	r := router.SetupRoutes(
		authHandler,
		laboratoryHandler,
		specificationHandler,
		categoryHandler,
		cfg.JWT.Secret,
		redisClient,
	)
	r.Use(middleware.ErrorHandler())
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
