package models

import (
	"time"
)

// User merepresentasikan entitas pengguna dalam sistem
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"` // Disembunyikan dari JSON response
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterInput format data JSON saat registrasi
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput format data JSON saat login
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
