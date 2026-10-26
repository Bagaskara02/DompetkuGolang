package main

import (
	"fmt"
	"log"

	"go_vibe/internal/config"
	"go_vibe/internal/database"
	"go_vibe/internal/handlers"
	"go_vibe/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Memuat konfigurasi environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[FATAL] Konfigurasi gagal dimuat: %v", err)
	}

	// 2. Inisialisasi koneksi database PostgreSQL & AutoMigrate GORM
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[FATAL] Database gagal diinisialisasi: %v", err)
	}
	log.Println("[INFO] Terhubung ke database PostgreSQL dan migrasi skema berhasil.")

	// 3. Inisialisasi Handler
	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	expenseHandler := handlers.NewExpenseHandler(db)

	// 4. Inisialisasi Router Gin
	r := gin.Default()

	// 5. Sajikan file antarmuka statis dari folder web/
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/index.html", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")
	r.StaticFile("/login.html", "./web/login.html")
	r.StaticFile("/register", "./web/register.html")
	r.StaticFile("/register.html", "./web/register.html")

	// 6. Rute REST API
	api := r.Group("/api")
	{
		// Rute Publik (Autentikasi)
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/logout", authHandler.Logout)

		// Rute Terproteksi (Wajib Login / Memeriksa JWT Cookie)
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			protected.GET("/me", authHandler.Me)

			// Catatan Pengeluaran
			protected.GET("/expenses", expenseHandler.GetExpenses)
			protected.GET("/expenses/summary", expenseHandler.GetSummary)
			protected.GET("/expenses/:id", expenseHandler.GetExpenseByID)
			protected.POST("/expenses", expenseHandler.CreateExpense)
			protected.PUT("/expenses/:id", expenseHandler.UpdateExpense)
			protected.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
		}
	}

	// 7. Menjalankan server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[INFO] Server DompetKu aktif di http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
