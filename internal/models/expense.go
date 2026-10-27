package models

import (
	"time"
)

// Expense merepresentasikan catatan transaksi keuangan pribadi (pengeluaran atau pemasukan)
type Expense struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Amount    float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	Type      string    `gorm:"size:20;not null;default:'pengeluaran';index" json:"type"` // "pemasukan" atau "pengeluaran"
	Category  string    `gorm:"size:100;not null;index" json:"category"`
	Date      time.Time `gorm:"type:date;not null" json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi ke User
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

// ExpenseInput format data JSON saat menambah atau mengubah catatan keuangan
type ExpenseInput struct {
	Title    string  `json:"title" binding:"required,min=1,max=255"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Type     string  `json:"type"` // "pemasukan" atau "pengeluaran" (default "pengeluaran")
	Category string  `json:"category" binding:"required,min=1,max=100"`
	Date     string  `json:"date" binding:"required"` // Format YYYY-MM-DD
}

// ExpenseSummary format ringkasan total saldo, pemasukan, dan pengeluaran
type ExpenseSummary struct {
	TotalIncome  float64 `json:"total_income"`  // Total Pemasukan
	TotalExpense float64 `json:"total_expense"` // Total Pengeluaran
	Balance      float64 `json:"balance"`       // Saldo (Total Pemasukan - Total Pengeluaran)
	TotalCount   int64   `json:"total_count"`
}
