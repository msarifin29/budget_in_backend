package model

import (
	"time"
)

type Income struct {
	Uid            string     `json:"uid" `
	Id             float64    `json:"id"`
	CategoryIncome string     `json:"category_income"` // Old entity , will be remove later
	TypeIncome     string     `json:"type_income"`
	Total          float64    `json:"total"`
	TransactionId  string     `json:"transaction_id"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	AccountId      string     `json:"account_id"`
	BankName       string     `json:"bank_name"`
	BankId         string     `json:"bank_id"`
	Cid            int32      `json:"c_id"`
}
type CreateIncomeRequest struct {
	Uid            string  `json:"uid" binding:"required"`
	CategoryIncome string  `json:"category_income"` // Old entity , will be remove later
	CategoryId     int32   `json:"category_id"`
	TypeIncome     string  `json:"type_income" binding:"required,type_income"`
	Total          float64 `json:"total" binding:"required,min=2000"`
	AccountId      string  `json:"account_id" binding:"required"`
	TransactionId  string  `json:"transaction_id" binding:"required"`
	CreatedAt      string  `json:"created_at"`
	BankName       string  `json:"bank_name"`
	BankId         string  `json:"bank_id"`
	Cid            int32   `json:"c_id"`
}
type CreateIncomeParams struct {
	Uid            string  `json:"uid" binding:"required"`
	CategoryIncome string  `json:"category_income"` // Old entity , will be remove later
	CategoryId     int32   `json:"category_id" binding:"required"`
	TypeIncome     string  `json:"type_income" binding:"required,type_income"`
	Total          float64 `json:"total" binding:"required,min=2000"`
	AccountId      string  `json:"account_id" binding:"required"`
	TransactionId  string  `json:"transaction_id"`
	CreatedAt      string  `json:"created_at"`
	BankName       string  `json:"bank_name"`
	BankId         string  `json:"bank_id"`
}
type GetIncomeParams struct {
	Uid            string `json:"uid" binding:"required"`
	CategoryIncome string `json:"category_income"` // Old entity , will be remove later
	CategoryId     int32  `json:"category_id"`
	TypeIncome     string `json:"type_income"`
	Limit          int32  `json:"limit"`
	Offset         int32  `json:"offset"`
}
type GetIncomeRequest struct {
	CategoryIncome string `form:"category_income"` // Old entity , will be remove later
	TypeIncome     string `form:"type_income"`
	CategoryId     int32  `form:"category_id"`
	Page           int32  `form:"page" binding:"required,min=1"`
	TotalPage      int32  `form:"total_page" binding:"required,min=5,max=10"`
}
type IncomeResponse struct {
	Uid            string          `json:"uid" `
	Id             float64         `json:"id"`
	CategoryIncome string          `json:"category_income"` // Old entity , will be remove later
	TypeIncome     string          `json:"type_income"`
	Total          float64         `json:"total"`
	TransactionId  string          `json:"transaction_id"`
	TCategory      CategoryReponse `json:"t_category"`
	CreatedAt      *time.Time      `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	AccountId      string          `json:"account_id"`
	BankName       string          `json:"bank_name"`
	BankId         string          `json:"bank_id"`
}
type IncomesResponse struct {
	Page      int32            `json:"page"`
	TotalPage int32            `json:"total_page"`
	LastPage  int32            `json:"last_page"`
	Total     int32            `json:"total"`
	Data      []IncomeResponse `json:"data"`
}
