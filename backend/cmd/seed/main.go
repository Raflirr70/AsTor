package seed

import (
	"astor/backend/entity"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []entity.Roles{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "PJ Gudang"},
		{ID: 3, Name: "PJ LABKOM 1"},
		{ID: 4, Name: "PJ LABKOM 2"},
		{ID: 5, Name: "PJ LABKOM 3"},
		{ID: 6, Name: "Asisten"},
	}
	for _, role := range roles {
		if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			continue
		}
	}
	log.Println("Seeder Role Success")
}
func SeedUsers(db *gorm.DB) {
	users := []entity.Users{
		{RoleID: 1, Firstname: "Super", Lastname: "Admin", Birthday: time.Now(), Address: "", Npm: "0000000000", Username: "superadmin", Password: "superadmin"},
		{RoleID: 2, Firstname: "Rafli", Lastname: "Radiansyah", Birthday: time.Now(), Address: "", Npm: "552012092", Username: "232401", Password: "232401"},
	}

	for _, user := range users {
		var exist entity.Users
		err := db.Where("username = ?", user.Username).First(&exist).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		// Hash password only for new users
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to seed user: %v", err)
		}
	}
	log.Println("Seeder User Success")
}
func SeedLocations(db *gorm.DB) {
	locations := []entity.Locations{
		{ID: 1, Name: "Labkom 1"},
		{ID: 2, Name: "Labkom 2"},
		{ID: 3, Name: "Labkom 3"},
		{ID: 4, Name: "Ruang Asisten"},
		{ID: 5, Name: "Gudang"},
	}

	for _, loc := range locations {
		var exist entity.Locations
		err := db.Where("name = ?", loc.Name).First(&exist).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err := db.Create(&loc).Error; err != nil {
			log.Fatalf("Failed to seed locations: %v", err)
		}
	}
	log.Println("Seeder Locations Success")
}
func SeedLabs(db *gorm.DB) {
	locations := []entity.Laboratories{
		{ID: 1, LocationID: 1, Capasity: 15},
		{ID: 2, LocationID: 2, Capasity: 22},
		{ID: 3, LocationID: 3, Capasity: 20},
	}

	for _, lab := range locations {
		var exist entity.Laboratories
		err := db.Where("location_id = ?", lab.LocationID).First(&exist).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err := db.Create(&lab).Error; err != nil {
			log.Fatalf("Failed to seed laboratories: %v", err)
		}
	}
	log.Println("Seeder Laboratories Success")
}
func SeedCategories(db *gorm.DB) {
	categories := []entity.Categories{
		{ID: 1, Name: "RAM", Code: "RAM"},
		{ID: 2, Name: "VGA", Code: "VGA"},
		{ID: 3, Name: "Motherboard", Code: "MOB"},
		{ID: 4, Name: "Power Supply", Code: "PSU"},
		{ID: 5, Name: "SSD", Code: "SSD"},
		{ID: 6, Name: "HDD", Code: "HDD"},
		{ID: 7, Name: "Processor", Code: "CPU"},
		{ID: 8, Name: "Keyboard", Code: "KEY"},
		{ID: 9, Name: "Monitor", Code: "MON"},
		{ID: 10, Name: "Mouse", Code: "MOS"},
		{ID: 11, Name: "Casing", Code: "CAS"},
		{ID: 12, Name: "Fan", Code: "FAN"},
		{ID: 13, Name: "Webcam", Code: "CAM"},
		{ID: 14, Name: "Headset", Code: "HDS"},
		{ID: 15, Name: "Speaker", Code: "SPK"},
		{ID: 16, Name: "Microphone", Code: "MIC"},
		{ID: 17, Name: "UPS", Code: "UPS"},
		{ID: 18, Name: "IPS", Code: "IPS"},
		{ID: 19, Name: "Printer", Code: "PRN"},
		{ID: 20, Name: "Scanner", Code: "SCN"},
		{ID: 21, Name: "Projector", Code: "PRJ"},
		{ID: 22, Name: "USB Hub", Code: "HUB"},
		{ID: 23, Name: "Adapter", Code: "ADP"},
		{ID: 24, Name: "Cable Power", Code: "CBP"},
		{ID: 25, Name: "Cable SATA", Code: "CBS"},
		{ID: 26, Name: "Cable VGA", Code: "CBV"},
		{ID: 27, Name: "Cable HDMI", Code: "CBH"},
		{ID: 28, Name: "Cable DVI", Code: "CBD"},
		{ID: 29, Name: "MousePad", Code: "MSP"},
		{ID: 30, Name: "Table", Code: "TBL"},
		{ID: 31, Name: "Chair", Code: "CHR"},
		{ID: 32, Name: "Switch", Code: "SWT"},
		{ID: 33, Name: "Mikrotik", Code: "MKR"},
		{ID: 34, Name: "Router", Code: "RTR"},
		{ID: 35, Name: "Modul", Code: "MDL"},
	}

	for _, category := range categories {
		var exist entity.Categories

		err := db.Where(
			"name = ? OR code = ? OR id = ?",
			category.Name, category.Code, category.ID,
		).First(&exist).Error

		if err == nil {
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to find category: %v", err)
		}

		if err := db.Create(&category).Error; err != nil {
			log.Fatalf("Failed to seed category %s: %v", category.Name, err)
		}
	}

	log.Println("Seeder Categories Success")
}
func SeedSpecifications(db *gorm.DB) {
	specifications := []entity.Specifications{
		{ID: 1, Name: "Storage", Unit: "GB"},
		{ID: 2, Name: "Capasity", Unit: "GB"},
		{ID: 3, Name: "Type", Unit: ""},
		{ID: 4, Name: "Transfer Speed", Unit: "MT/s"},
		{ID: 5, Name: "Read Speed", Unit: "MB/s"},
		{ID: 6, Name: "Core", Unit: "Core"},
		{ID: 7, Name: "Thread", Unit: "Thread"},
		{ID: 8, Name: "Power", Unit: "Watt"},
		{ID: 9, Name: "Screen Size", Unit: "Inch"},
		{ID: 10, Name: "Quantity", Unit: "Pcs"},
		{ID: 11, Name: "Clock Speed", Unit: "MHz"},
		{ID: 12, Name: "Port", Unit: "Port"},
	}

	for _, spec := range specifications {
		var exist entity.Specifications
		err := db.Where("name = ?", spec.Name).First(&exist).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err := db.Create(&spec).Error; err != nil {
			log.Fatalf("Failed to seed Specifications: %v", err)
		}
	}
	log.Println("Seeder Specifications Success")
}
func SeedCategorySpecifications(db *gorm.DB) {
	css := []entity.CategorySpecifications{
		{ID: 1, CategoryID: 1, SpecificationID: 2}, {ID: 2, CategoryID: 1, SpecificationID: 3}, {ID: 3, CategoryID: 1, SpecificationID: 4},
		{ID: 4, CategoryID: 2, SpecificationID: 2}, {ID: 5, CategoryID: 2, SpecificationID: 3}, {ID: 6, CategoryID: 2, SpecificationID: 4},
		{ID: 7, CategoryID: 3, SpecificationID: 3}, {ID: 8, CategoryID: 3, SpecificationID: 3},
		{ID: 9, CategoryID: 4, SpecificationID: 8},
		{ID: 10, CategoryID: 5, SpecificationID: 1}, {ID: 11, CategoryID: 5, SpecificationID: 5},
		{ID: 12, CategoryID: 6, SpecificationID: 1}, {ID: 13, CategoryID: 6, SpecificationID: 5},
		{ID: 14, CategoryID: 7, SpecificationID: 11}, {ID: 15, CategoryID: 7, SpecificationID: 6}, {ID: 16, CategoryID: 7, SpecificationID: 7},
		{ID: 17, CategoryID: 9, SpecificationID: 8}, {ID: 18, CategoryID: 9, SpecificationID: 9},
		{ID: 19, CategoryID: 24, SpecificationID: 10},
		{ID: 20, CategoryID: 25, SpecificationID: 10},
		{ID: 21, CategoryID: 26, SpecificationID: 10},
		{ID: 22, CategoryID: 27, SpecificationID: 10},
		{ID: 23, CategoryID: 28, SpecificationID: 10},
		{ID: 24, CategoryID: 29, SpecificationID: 10},
		{ID: 25, CategoryID: 32, SpecificationID: 12},
		{ID: 26, CategoryID: 35, SpecificationID: 10},
	}

	for _, cs := range css {
		if err := db.Create(&cs).Error; err != nil {
			log.Fatalf("Failed to seed Specifications: %v", err)
		}
	}
	log.Println("Seeder Specifications Success")
}
