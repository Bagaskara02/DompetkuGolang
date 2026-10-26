package handlers

import (
	"net/http"
	"strings"

	"go_vibe/internal/middleware"
	"go_vibe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthHandler(db *gorm.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Register menangani pembuatan akun baru
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data input tidak valid. Pastikan username minimal 3 karakter, email valid, dan password minimal 6 karakter.",
		})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	// Periksa apakah username atau email sudah terdaftar
	var existing models.User
	if err := h.db.Where("LOWER(username) = ? OR LOWER(email) = ?", strings.ToLower(input.Username), input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username atau email sudah digunakan oleh akun lain.",
		})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := middleware.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Terjadi kesalahan saat memproses password.",
		})
		return
	}

	newUser := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan akun baru ke database.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi akun berhasil. Silakan login.",
		"user": gin.H{
			"id":       newUser.ID,
			"username": newUser.Username,
			"email":    newUser.Email,
		},
	})
}

// Login memverifikasi kredensial dan membuat HTTP-only cookie berisi JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username dan password wajib diisi.",
		})
		return
	}

	input.Username = strings.TrimSpace(input.Username)

	var user models.User
	if err := h.db.Where("LOWER(username) = ?", strings.ToLower(input.Username)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username atau password salah.",
		})
		return
	}

	// Periksa kecocokan password hash
	if !middleware.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username atau password salah.",
		})
		return
	}

	// Buat token JWT
	token, err := middleware.GenerateToken(user.ID, user.Username, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat sesi login.",
		})
		return
	}

	// Simpan token di HTTP-only cookie (Aman dari pembacaan JavaScript/XSS)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		middleware.CookieName,
		token,
		86400, // 24 jam dalam detik
		"/",
		"",
		false, // Set true jika menggunakan protokol HTTPS
		true,  // HttpOnly = true
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil.",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout membersihkan HTTP-only cookie
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		middleware.CookieName,
		"",
		-1, // Langsung kedaluwarsa
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil.",
	})
}

// Me mengembalikan profil pengguna yang sedang login
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	userID := userIDVal.(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}
