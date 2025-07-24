package database

import (
	"app/config"
	"app/model"
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB initializes the database connection
func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Invalid port number: %s", p))
	}

	dsn := fmt.Sprintf(
		"host=db port=%s password=%s user=%s dbname=%s sslmode=disable",
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	fmt.Println("Connected to database successfully")
	DB.AutoMigrate(
		&model.Category{},
		&model.Comment{},
		&model.Item{},
		&model.Like{},
		&model.Order{},
		&model.Post{},
		&model.Review{},
		&model.User{},
	)
	fmt.Println("Database migrated successfully")
}
