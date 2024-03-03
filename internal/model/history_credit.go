package model

import "time"

type HistoryCredit struct {
	CreditId    float64   `json:"credit_id"`
	Id          float64   `json:"id"`
	Th          float64   `json:"th"`
	Total       float64   `json:"total"`
	Status      string    `json:"status"`
	TypePayment string    `json:"type_payment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PaymentTime int       `json:"payment_time"`
}
type CreateHistoryCredit struct {
	CreditId    float64 `json:"credit_id" binding:"required"`
	Th          float64 `json:"th" binding:"required,min=1"`
	Total       float64 `json:"total" binding:"required,min=2000"`
	Status      string  `json:"status" binding:"required,status_credit"`
	TypePayment string  `json:"type_payment"`
	PaymentTime int     `json:"payment_time"`
}
type UpdateHistoryCreditRequest struct {
	Uid         string  `json:"uid" binding:"required"`
	Id          float64 `json:"id" binding:"required"`
	Status      string  `json:"status" binding:"required,status_credit"`
	TypePayment string  `json:"type_payment" binding:"required,expense_type"`
}

func NewHistoryCredit(history HistoryCredit) *HistoryCredit {
	return &HistoryCredit{
		CreditId:    history.CreditId,
		Id:          history.Id,
		Th:          history.Th,
		Total:       history.Total,
		Status:      history.Status,
		TypePayment: history.TypePayment,
		PaymentTime: history.PaymentTime,
		CreatedAt:   history.CreatedAt,
		UpdatedAt:   history.UpdatedAt,
	}
}

type HistoriesCreditsResponse struct {
	Page      int32           `json:"page"`
	TotalPage int32           `json:"total_page"`
	LastPage  int32           `json:"last_page"`
	Total     int32           `json:"total"`
	Data      []HistoryCredit `json:"data"`
}
