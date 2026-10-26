package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config menyimpan seluruh konfigurasi aplikasi dari environment variable
type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
}

// LoadConfig membaca environment variable dari file .env (jika ada) dan dari sistem
func LoadConfig() (*Config, error) {
	// Membaca file .env jika ada (tidak wajib ada jika dijalankan di container/production)
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("konfigurasi wajib DATABASE_URL belum diatur di file .env atau environment variable")
	}

	// Dukung JWT_SECRET atau SESSION_SECRET sesuai spesifikasi
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SESSION_SECRET")
	}
	if jwtSecret == "" {
		return nil, errors.New("konfigurasi wajib JWT_SECRET / SESSION_SECRET belum diatur di file .env atau environment variable")
	}

	return &Config{
		AppPort:     port,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}, nil
}
