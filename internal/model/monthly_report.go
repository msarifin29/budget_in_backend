package model

import (
	"database/sql"
	"time"
)

type MonthlyParams struct {
	Uid   string `json:"uid" `
	Year  string `json:"year" `
	Month string `json:"month" `
}

type MonthlyRequest struct {
	Year  string `form:"year" binding:"required"`
	Month string `form:"month" binding:"required"`
}
type ParamMonthlyReport struct {
	Uid string `uri:"uid" json:"uid" binding:"required"`
}
type MonthlyResponse struct {
	Date          *time.Time `json:"date"`
	TotalIncomes  float64    `json:"total_incomes"`
	TotalExpenses float64    `json:"total_expenses"`
}
type MonthlyReport struct {
	Uid          sql.NullString  `json:"uid"`
	Month        sql.NullFloat64 `json:"month"`
	Year         sql.NullFloat64 `json:"year"`
	TotalIncome  sql.NullFloat64 `json:"total_income"`
	TotalExpense sql.NullFloat64 `json:"total_expense"`
}
type MonthlyReportResponse struct {
	Month        float64 `json:"month"`
	Year         float64 `json:"year"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

// Details
type ParamMonthlyReportDetail struct {
	Month string `json:"month" binding:"required"`
	Uid   string `json:"uid" binding:"required"`
}
type RequestMonthlyReportDetail struct {
	Month string `form:"month" binding:"required"`
}
type MonthlyXDetail struct {
	Month         string          `json:"month"`
	Uid           string          `json:"uid"`
	Id            float64         `json:"id"`
	ExpenseType   string          `json:"expense_type"`
	Total         float64         `json:"total"`
	Status        string          `json:"status"`
	Notes         string          `json:"notes"`
	TransactionId string          `json:"transaction_id"`
	TCategory     CategoryReponse `json:"t_category"`
	CreatedAt     *time.Time      `json:"created_at"`
}
type MonthlyIDetail struct {
	Month         string          `json:"month"`
	Uid           string          `json:"uid" `
	Id            float64         `json:"id"`
	TypeIncome    string          `json:"type_income"`
	Total         float64         `json:"total"`
	TransactionId string          `json:"transaction_id"`
	TCategory     CategoryReponse `json:"t_category"`
	CreatedAt     *time.Time      `json:"created_at"`
}
type MonthlyReportDetailResponse struct {
	ExpensesRecords []MonthlyXDetail `json:"expenses"`
	IncomesRecords  []MonthlyIDetail `json:"incomes"`
}

type MonthlyReportCategoryExpense struct {
	CategoryId sql.NullString  `json:"category_id"`
	Title      sql.NullString  `json:"title"`
	Total      sql.NullFloat64 `json:"total"`
}

type MonthlyReportCategoryExpenseResponse struct {
	CategoryId string  `json:"category_id"`
	Title      string  `json:"title"`
	Total      float64 `json:"total"`
}
type MonthlyReportCategoryIncome struct {
	CategoryId sql.NullString  `json:"category_id"`
	Title      sql.NullString  `json:"title"`
	Total      sql.NullFloat64 `json:"total"`
}

type MonthlyReportCategoryIncomeResponse struct {
	CategoryId string  `json:"category_id"`
	Title      string  `json:"title"`
	Total      float64 `json:"total"`
}

type MonthlyReportCategoryResponse struct {
	Incomes  []MonthlyReportCategoryIncomeResponse  `json:"incomes"`
	Expenses []MonthlyReportCategoryExpenseResponse `json:"expenses"`
}

type ParamMonthlyReportCategory struct {
	Uid   string `json:"uid" binding:"required"`
	Month string `json:"month" binding:"required"`
}
type RequestMonthlyReportCategory struct {
	Month string `form:"month" binding:"required"`
}
