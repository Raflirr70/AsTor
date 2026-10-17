package main

import (
	"astor/backend/config"
	"astor/backend/entity"
	"astor/backend/pkg/database"
	"astor/backend/pkg/logger"
	"log"

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

	// 3. Migration
	if err := db.AutoMigrate(
		&entity.Roles{},
		&entity.Users{},
		&entity.ActivityLogs{},
		&entity.Notifications{},
		&entity.Loans{},
		&entity.Locations{},
		&entity.Categories{},
		&entity.Specifications{},
		&entity.Assets{},
		&entity.AssetSpecifications{},
		&entity.Computers{},
		&entity.ComputerComponents{},
		&entity.Laboratories{},
		&entity.ComputerLaboratories{},
		&entity.AssetPlacements{},
		&entity.LoanDetails{},
		&entity.Reports{},
		&entity.TranferRequests{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	logger.Info.Println("Migration Success")

	// seed.SeedRoles(db)
}
