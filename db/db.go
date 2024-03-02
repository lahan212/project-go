package db

import (
	"project-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func Init() {
	// Connect to MySQL database
	dsn := "root:@tcp(localhost:3306)/golang_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// Auto Migrate the User model
	DB.AutoMigrate(&models.Users{})
}
