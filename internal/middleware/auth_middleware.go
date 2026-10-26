package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const CookieName = "dompetku_token"

// JWTClaims merepresentasikan data payload dalam token JWT
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// HashPassword membuat hash dari password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash memverifikasi apakah password sesuai dengan hash yang tersimpan
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken membuat token JWT baru dengan masa berlaku 24 jam
func GenerateToken(userID uint, username, secret string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken memvalidasi string token JWT dan mengembalikan claims
func ValidateToken(tokenStr, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid atau sudah kedaluwarsa")
}

// AuthRequired adalah middleware Gin untuk memverifikasi cookie JWT
func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(CookieName)
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Akses ditolak. Silakan login terlebih dahulu.",
			})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenStr, secret)
		if err != nil {
			// Jika token tidak valid, bersihkan cookie lama
			c.SetCookie(CookieName, "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Sesi Anda telah kedaluwarsa. Silakan login kembali.",
			})
			c.Abort()
			return
		}

		// Simpan informasi user di konteks Gin untuk digunakan di handler berikutnya
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
