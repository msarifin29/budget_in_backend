package model

import "time"

type MonthlyParams struct {
	Uid   string `json:"uid" `
	Year  string `json:"year" `
	Month string `json:"month" `
}

type MonthlyRequest struct {
	Year  string `form:"year" binding:"required"`
	Month string `form:"month" binding:"required"`
}

type MonthlyResponse struct {
	Date          time.Time `json:"date"`
	TotalIncomes  float64   `json:"total_incomes"`
	TotalExpenses float64   `json:"total_expenses"`
}
