package model

import (
	"time"
)

type Expense struct {
	Uid           string     `json:"uid" binding:"required"`
	Id            float64    `json:"id"`
	ExpenseType   string     `json:"expense_type"`
	Total         float64    `json:"total"`
	Category      string     `json:"category"` // Not used will be remove later
	Status        string     `json:"status"`
	Notes         string     `json:"notes"`
	TransactionId string     `json:"transaction_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	AccountId     string     `json:"account_id"`
	BankName      string     `json:"bank_name"`
	BankId        string     `json:"bank_id"`
	Cid           float64    `json:"c_id"`
}

type CreateExpenseRequest struct {
	Uid           string  `json:"uid" binding:"required"`
	ExpenseType   string  `json:"expense_type" binding:"required,expense_type"`
	Category      string  `json:"category"` // Not used will be remove later
	CategoryId    float64 `json:"category_id" binding:"required"`
	Total         float64 `json:"total" binding:"required,min=2000"`
	Notes         string  `json:"notes"`
	AccountId     string  `json:"account_id" binding:"required"`
	TransactionId string  `json:"transaction_id"`
	CreatedAt     string  `json:"created_at"`
	BankName      string  `json:"bank_name"`
	BankId        string  `json:"bank_id"`
}
type CreateExpenseParams struct {
	Uid         string  `json:"uid" binding:"required"`
	ExpenseType string  `json:"expense_type" binding:"required,expense_type"`
	Category    string  `json:"category"` // Not used will be remove later
	CategoryId  float64 `json:"category_id" binding:"required"`
	Total       float64 `json:"total" binding:"required,min=2000"`
	Notes       string  `json:"notes"`
	AccountId   string  `json:"account_id" binding:"required"`
	CreatedAt   string  `json:"created_at"`
	BankName    string  `json:"bank_name"`
	BankId      string  `json:"bank_id"`
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
	Category    string `json:"category"` // Not used
	Id          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

type ExpenseParamWithId struct {
	Id float64 `uri:"id" json:"id" binding:"required"`
}

type GetExpenseRequest struct {
	Status      string `form:"status"`
	ExpenseType string `form:"expense_type"`
	Category    string `form:"category"` // Not Used
	Id          int    `form:"id"`
	CreatedAt   string `form:"created_at"`
	Page        int32  `form:"page" binding:"required,min=1"`
	TotalPage   int32  `form:"total_page" binding:"required,min=5,max=10"`
}

type ExpenseResponse struct {
	Uid           string          `json:"uid" binding:"required"`
	Id            float64         `json:"id"`
	ExpenseType   string          `json:"expense_type"`
	Total         float64         `json:"total"`
	Category      string          `json:"category"`
	Status        string          `json:"status"`
	Notes         string          `json:"notes"`
	TransactionId string          `json:"transaction_id"`
	TCategory     CategoryReponse `json:"t_category"`
	CreatedAt     *time.Time      `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
	AccountId     string          `json:"account_id"`
	BankName      string          `json:"bank_name"`
	BankId        string          `json:"bank_id"`
}

type ExpensesResponse struct {
	Page      int32             `json:"page"`
	TotalPage int32             `json:"total_page"`
	LastPage  int32             `json:"last_page"`
	Total     int32             `json:"total"`
	Data      []ExpenseResponse `json:"data"`
}
