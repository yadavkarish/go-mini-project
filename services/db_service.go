package services

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitializeDatabase initializes the database connection
func InitializeDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logs.Fatalf("Failed to configure database: %v", err)
	}

	sqlDB.SetMaxOpenConns(20) // Adjust based on system capacity
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	return db
}
