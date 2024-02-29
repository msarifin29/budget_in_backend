package model

import (
	"time"
)

type Expense struct {
	Uid         string    `json:"uid" binding:"required"`
	Id          float64   `json:"id"`
	ExpenseType string    `json:"expense_type"`
	Total       float64   `json:"total"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Uid         string  `json:"uid" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required,expense_type"`
	Total       float64 `json:"total" binding:"required,min=2000"`
	Notes       string  `json:"notes"`
}

type UpdateExpenseRequest struct {
	Id          float64 `json:"id" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required,expense_type"`
	Total       float64 `json:"total" binding:"required,min=2000"`
	Notes       string  `json:"notes"`
}

type GetExpenseParams struct {
	Uid    string `json:"uid" binding:"required"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ExpenseParamWithId struct {
	Id float64 `uri:"id" json:"id" binding:"required"`
}

type GetExpenseRequest struct {
	// Uid       string `form:"uid" binding:"required"`
	Page      int32 `form:"page" binding:"required,min=1"`
	TotalPage int32 `form:"total_page" binding:"required,min=5"`
}

type ExpenseResponse struct {
	Uid         string    `json:"uid" binding:"required"`
	Id          float64   `json:"id"`
	ExpenseType string    `json:"expense_type"`
	Total       float64   `json:"total"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpensesResponse struct {
	Page      int32     `json:"page"`
	TotalPage int32     `json:"total_page"`
	Data      []Expense `json:"data"`
}
