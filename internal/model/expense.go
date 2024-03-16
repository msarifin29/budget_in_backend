package model

import (
	"time"
)

type Expense struct {
	Uid         string     `json:"uid" binding:"required"`
	Id          float64    `json:"id"`
	ExpenseType string     `json:"expense_type"`
	Total       float64    `json:"total"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Notes       string     `json:"notes"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Uid         string  `json:"uid" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required,expense_type"`
	Category    string  `json:"category" binding:"required,category_expense"`
	Total       float64 `json:"total" binding:"required,min=2000"`
	Notes       string  `json:"notes"`
	AccountId   string  `json:"account_id" binding:"required"`
}

type UpdateExpenseRequest struct {
	Id          float64 `json:"id" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required,expense_type"`
	AccountId   string  `json:"account_id" binding:"required"`
}

type GetExpenseParams struct {
	Uid         string `json:"uid" binding:"required"`
	Status      string `json:"status"`
	ExpenseType string `json:"expense_type"`
	Category    string `json:"category"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type ExpenseParamWithId struct {
	Id float64 `uri:"id" json:"id" binding:"required"`
}

type GetExpenseRequest struct {
	Status      string `form:"status" binding:"required,status"`
	ExpenseType string `form:"expense_type"`
	Category    string `form:"category"`
	Page        int32  `form:"page" binding:"required,min=1"`
	TotalPage   int32  `form:"total_page" binding:"required,min=5,max=10"`
}

type ExpenseResponse struct {
	Uid         string     `json:"uid" binding:"required"`
	Id          float64    `json:"id"`
	ExpenseType string     `json:"expense_type"`
	Total       float64    `json:"total"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Notes       string     `json:"notes"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ExpensesResponse struct {
	Page      int32     `json:"page"`
	TotalPage int32     `json:"total_page"`
	LastPage  int32     `json:"last_page"`
	Total     int32     `json:"total"`
	Data      []Expense `json:"data"`
}
