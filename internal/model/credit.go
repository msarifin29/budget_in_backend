package model

import "time"

type Credit struct {
	Uid            string     `json:"uid"`
	Id             float64    `json:"id"`
	CategoryCredit string     `json:"category_credit"`
	TypeCredit     string     `json:"type_credit"`
	Total          float64    `json:"total"`
	LoanTerm       float64    `json:"loan_term"`
	PaymentTime    int        `json:"payment_time"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	StatusCredit   string     `json:"status_credit"`
	Installment    float64    `json:"installment"`
}
type CreditResponse struct {
	Uid            string          `json:"uid"`
	Id             float64         `json:"id"`
	CategoryCredit string          `json:"category_credit"`
	TypeCredit     string          `json:"type_credit"`
	Total          float64         `json:"total"`
	LoanTerm       float64         `json:"loan_term"`
	PaymentTime    int             `json:"payment_time"`
	CreatedAt      *time.Time      `json:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at"`
	StatusCredit   string          `json:"status_credit"`
	Installment    float64         `json:"installment"`
	TCategory      CategoryReponse `json:"t_category"`
}
type CreateCreditRequest struct {
	Uid            string  `json:"uid" binding:"required"`
	CategoryCredit string  `json:"category_credit"`
	CategoryId     int32   `json:"category_id" binding:"required"`
	TypeCredit     string  `json:"type_credit" binding:"required,type_credit"`
	LoanTerm       float64 `json:"loan_term" binding:"required,min=1"`
	Installment    float64 `json:"installment" binding:"required,min=2000"`
	PaymentTime    int     `json:"payment_time" binding:"required,min=1,max=31"`
}
type UpdateCreditRequest struct {
	Uid          string  `json:"uid" binding:"required"`
	Id           float64 `json:"id" binding:"required"`
	StatusCredit string  `json:"status_credit" binding:"required,status_credit"`
}
type GetCreditRequest struct {
	Uid string  `json:"uid" binding:"required"`
	Id  float64 `json:"id" binding:"required"`
}
type GetCreditParams struct {
	Uid    string `json:"uid" binding:"required"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}
type GetCreditsRequest struct {
	Page      int32 `form:"page" binding:"required,min=1"`
	TotalPage int32 `form:"total_page" binding:"required,min=5,max=10"`
}

func NewCredit(credit Credit) *Credit {
	return &Credit{
		Uid:            credit.Uid,
		Id:             credit.Id,
		CategoryCredit: credit.CategoryCredit,
		TypeCredit:     credit.TypeCredit,
		Total:          credit.Total,
		LoanTerm:       credit.LoanTerm,
		StatusCredit:   credit.StatusCredit,
		Installment:    credit.Installment,
		PaymentTime:    credit.PaymentTime,
		CreatedAt:      credit.CreatedAt,
		UpdatedAt:      credit.UpdatedAt,
	}
}

type CreditsResponse struct {
	Page      int32            `json:"page"`
	TotalPage int32            `json:"total_page"`
	LastPage  int32            `json:"last_page"`
	Total     int32            `json:"total"`
	Data      []CreditResponse `json:"data"`
}
