package config

import (
	"os"
	"testing"
)

func TestLoadConfigSuccess(t *testing.T) {
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("JWT_SECRET", "super-secret-key")
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig gagal: %v", err)
	}

	if cfg.AppPort != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.AppPort)
	}
	if cfg.DatabaseURL != "postgres://test:test@localhost:5432/testdb" {
		t.Errorf("Expected database URL, got %s", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "super-secret-key" {
		t.Errorf("Expected JWT secret, got %s", cfg.JWTSecret)
	}
}

func TestLoadConfigMissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	os.Setenv("JWT_SECRET", "secret")
	defer os.Unsetenv("JWT_SECRET")

	_, err := LoadConfig()
	if err == nil {
		t.Errorf("LoadConfig seharusnya menghasilkan error jika DATABASE_URL tidak diset")
	}
}
