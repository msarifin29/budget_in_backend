package model

import "time"

type Income struct {
	Uid            string     `json:"uid" `
	Id             float64    `json:"id"`
	CategoryIncome string     `json:"category_income"`
	TypeIncome     string     `json:"type_income"`
	Total          float64    `json:"total"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
type CreateIncomeRequest struct {
	Uid            string  `json:"uid" binding:"required"`
	CategoryIncome string  `json:"category_income" binding:"required,category_income"`
	TypeIncome     string  `json:"type_income" binding:"required,type_income"`
	Total          float64 `json:"total" binding:"required,min=2000"`
	AccountId      string  `json:"account_id" binding:"required"`
}
type GetIncomeParams struct {
	Uid            string `json:"uid" binding:"required"`
	CategoryIncome string `json:"category_income"`
	TypeIncome     string `json:"type_income"`
	Limit          int32  `json:"limit"`
	Offset         int32  `json:"offset"`
}
type GetIncomeRequest struct {
	CategoryIncome string `form:"category_income"`
	TypeIncome     string `form:"type_income"`
	Page           int32  `form:"page" binding:"required,min=1"`
	TotalPage      int32  `form:"total_page" binding:"required,min=5,max=10"`
}
type IncomeResponse struct {
	Uid            string     `json:"uid" `
	Id             float64    `json:"id"`
	CategoryIncome string     `json:"category_income"`
	TypeIncome     string     `json:"type_income"`
	Total          float64    `json:"total"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
type IncomesResponse struct {
	Page      int32    `json:"page"`
	TotalPage int32    `json:"total_page"`
	LastPage  int32    `json:"last_page"`
	Total     int32    `json:"total"`
	Data      []Income `json:"data"`
}
