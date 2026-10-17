package main

import (
	"context"
	"log"

	"astor/backend/config"
	"astor/backend/pkg/database"
	"astor/backend/pkg/logger"
	"astor/backend/pkg/ngrok"
	router "astor/backend/route"

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

	// 3. Setup Dependency
	//repo
	//usecase
	//handler

	// 3. Setup Router
	r := router.SetupRoutes()
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
