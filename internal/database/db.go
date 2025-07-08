package database

import (
	"duit-pasutri-be/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(models.User{}, models.Account{}, models.Category{}, models.Transaction{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
