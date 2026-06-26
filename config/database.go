package config

import (
	"fmt"
	"log"
	"os"

	"manifest-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.ShippingAgent{},
		&models.Vessel{},
		&models.Manifest{},
		&models.ManifestDetail{},
		&models.BC11{},
		&models.NPE{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	return db
}
