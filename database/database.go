package database

import (
	"fmt"
	"gorm-rdbms/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global variable
var DB *gorm.DB

func Database()  {
	// Connect to database
	dsn := "root:@tcp(127.0.0.1:3306)/gorm-rdbms?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Connection test
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to database")

	// Migrate the schema
	errMigrate := DB.AutoMigrate(models.Post{}, models.Author{}, models.Category{}, models.News{}, models.Articles{}, models.Author{})
	if errMigrate != nil {
		panic("failed to migrate")
	}
	fmt.Println("Migrate success!")
}
