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
