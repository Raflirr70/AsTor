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
			log.Fatalf("Fail Run Seed Role : %v", err)
		}
	}
	log.Fatalf("Seeder Role Success")
}
func SeedUsers(db *gorm.DB) {
	users := []entity.Users{
		{RoleID: 1, Firstname: "Super", Lastname: "Admin", Birthday: time.Now(), Address: "", Npm: "0000000000", Username: "superadmin", Password: "superadmin"},
		{RoleID: 2, Firstname: "Ujang", Lastname: "Slipnot", Birthday: time.Now(), Address: "", Npm: "0000000001", Username: "12345", Password: "12345"},
	}

	for _, user := range users {
		var exist entity.Users
		err := db.Where("username = ?", user.Username).First(&exist).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to find user: %v", err)
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
