package database

import (
	"errors"
	"fmt"
	"time"

	"go_vibe/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB menginisialisasi koneksi GORM ke PostgreSQL dan menjalankan AutoMigrate
func InitDB(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL tidak boleh kosong")
	}

	// Buka koneksi GORM
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database PostgreSQL: %v", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengonfigurasi pool database: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrate tabel User dan Expense
	err = db.AutoMigrate(&models.User{}, &models.Expense{})
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi skema database: %v", err)
	}

	return db, nil
}
