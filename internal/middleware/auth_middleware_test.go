package middleware

import (
	"testing"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "rahasia123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword gagal: %v", err)
	}

	if hash == password {
		t.Errorf("Password tidak boleh sama persis dengan hash")
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("CheckPasswordHash harus bernilai true untuk password yang cocok")
	}

	if CheckPasswordHash("passwordsalah", hash) {
		t.Errorf("CheckPasswordHash harus bernilai false untuk password yang salah")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "kunci-rahasia-test-12345"
	userID := uint(10)
	username := "budi"

	token, err := GenerateToken(userID, username, secret)
	if err != nil {
		t.Fatalf("GenerateToken gagal: %v", err)
	}

	if token == "" {
		t.Fatalf("Token tidak boleh kosong")
	}

	// Validasi dengan secret yang benar
	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken gagal: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}

	// Validasi dengan secret yang salah
	_, err = ValidateToken(token, "secret-yang-salah")
	if err == nil {
		t.Errorf("Validasi token dengan secret yang salah seharusnya menghasilkan error")
	}
}
