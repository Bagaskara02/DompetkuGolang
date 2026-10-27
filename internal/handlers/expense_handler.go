package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_vibe/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	db *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{
		db: db,
	}
}

// helper untuk mengambil userID dari context Gin
func getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := val.(uint)
	return userID, ok
}

// GetExpenses mengambil seluruh daftar transaksi milik user yang sedang login
// Mendukung filter query: ?search=xxx & category=xxx & type=xxx (pemasukan/pengeluaran)
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	category := strings.TrimSpace(c.Query("category"))
	txType := strings.TrimSpace(strings.ToLower(c.Query("type")))

	query := h.db.Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if category != "" && category != "Semua" {
		query = query.Where("category = ?", category)
	}

	if txType != "" && txType != "semua" {
		query = query.Where("type = ?", txType)
	}

	var expenses []models.Expense
	if err := query.Order("date DESC, id DESC").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
	})
}

// GetSummary menghitung ringkasan total pemasukan, total pengeluaran, saldo, dan jumlah transaksi
func (h *ExpenseHandler) GetSummary(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	category := strings.TrimSpace(c.Query("category"))

	query := h.db.Model(&models.Expense{}).Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if category != "" && category != "Semua" {
		query = query.Where("category = ?", category)
	}

	type SummaryResult struct {
		TotalIncome  float64 `gorm:"column:total_income"`
		TotalExpense float64 `gorm:"column:total_expense"`
		TotalCount   int64   `gorm:"column:total_count"`
	}

	var res SummaryResult
	err := query.Select(`
		COALESCE(SUM(CASE WHEN type = 'pemasukan' THEN amount ELSE 0 END), 0) as total_income,
		COALESCE(SUM(CASE WHEN type != 'pemasukan' THEN amount ELSE 0 END), 0) as total_expense,
		COUNT(*) as total_count
	`).Scan(&res).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung ringkasan saldo."})
		return
	}

	balance := res.TotalIncome - res.TotalExpense

	c.JSON(http.StatusOK, models.ExpenseSummary{
		TotalIncome:  res.TotalIncome,
		TotalExpense: res.TotalExpense,
		Balance:      balance,
		TotalCount:   res.TotalCount,
	})
}

// GetExpenseByID mengambil detail satu catatan pengeluaran berdasarkan ID
func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengeluaran tidak valid."})
		return
	}

	var expense models.Expense
	// Batasan ketat: hanya bisa mengambil data jika ID dan user_id cocok
	if err := h.db.Where("id = ? AND user_id = ?", uint(id), userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan pengeluaran tidak ditemukan atau Anda tidak memiliki akses."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": expense})
}

// CreateExpense menambahkan catatan pengeluaran baru
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	var input models.ExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua kolom (nama, nominal, kategori, tanggal) wajib diisi dengan benar."})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	txType := strings.TrimSpace(strings.ToLower(input.Type))
	if txType != "pemasukan" {
		txType = "pengeluaran"
	}

	newExpense := models.Expense{
		UserID:   userID,
		Title:    strings.TrimSpace(input.Title),
		Amount:   input.Amount,
		Type:     txType,
		Category: strings.TrimSpace(input.Category),
		Date:     parsedDate,
	}

	if err := h.db.Create(&newExpense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Catatan transaksi berhasil ditambahkan.",
		"expense": newExpense,
	})
}

// UpdateExpense memperbarui data catatan pengeluaran / pemasukan
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID transaksi tidak valid."})
		return
	}

	var expense models.Expense
	// Pastikan data milik user yang sedang login
	if err := h.db.Where("id = ? AND user_id = ?", uint(id), userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan transaksi tidak ditemukan atau Anda tidak memiliki akses."})
		return
	}

	var input models.ExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua kolom wajib diisi dengan benar."})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	if input.Type != "" {
		txType := strings.TrimSpace(strings.ToLower(input.Type))
		if txType == "pemasukan" || txType == "pengeluaran" {
			expense.Type = txType
		}
	}

	expense.Title = strings.TrimSpace(input.Title)
	expense.Amount = input.Amount
	expense.Category = strings.TrimSpace(input.Category)
	expense.Date = parsedDate

	if err := h.db.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui catatan transaksi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan transaksi berhasil diperbarui.",
		"expense": expense,
	})
}

// DeleteExpense menghapus catatan pengeluaran
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengeluaran tidak valid."})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", uint(id), userID).Delete(&models.Expense{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus catatan pengeluaran."})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan pengeluaran tidak ditemukan atau Anda tidak memiliki akses."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan pengeluaran berhasil dihapus.",
	})
}
