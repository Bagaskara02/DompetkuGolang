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

// GetExpenses mengambil seluruh daftar pengeluaran milik user yang sedang login
// Mendukung filter query: ?search=xxx & category=xxx
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi."})
		return
	}

	search := strings.TrimSpace(c.Query("search"))
	category := strings.TrimSpace(c.Query("category"))

	query := h.db.Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if category != "" && category != "Semua" {
		query = query.Where("category = ?", category)
	}

	var expenses []models.Expense
	if err := query.Order("date DESC, id DESC").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengeluaran."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
	})
}

// GetSummary menghitung ringkasan total nominal dan jumlah catatan pengeluaran user
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

	var totalCount int64
	query.Count(&totalCount)

	type SumResult struct {
		Total float64
	}
	var sumResult SumResult
	query.Select("COALESCE(SUM(amount), 0) as total").Scan(&sumResult)

	c.JSON(http.StatusOK, models.ExpenseSummary{
		TotalAmount: sumResult.Total,
		TotalCount:  totalCount,
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

	newExpense := models.Expense{
		UserID:   userID,
		Title:    strings.TrimSpace(input.Title),
		Amount:   input.Amount,
		Category: strings.TrimSpace(input.Category),
		Date:     parsedDate,
	}

	if err := h.db.Create(&newExpense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengeluaran."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Catatan pengeluaran berhasil ditambahkan.",
		"expense": newExpense,
	})
}

// UpdateExpense memperbarui data catatan pengeluaran
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
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
	// Pastikan data milik user yang sedang login
	if err := h.db.Where("id = ? AND user_id = ?", uint(id), userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Catatan pengeluaran tidak ditemukan atau Anda tidak memiliki akses."})
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

	expense.Title = strings.TrimSpace(input.Title)
	expense.Amount = input.Amount
	expense.Category = strings.TrimSpace(input.Category)
	expense.Date = parsedDate

	if err := h.db.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui catatan pengeluaran."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Catatan pengeluaran berhasil diperbarui.",
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
