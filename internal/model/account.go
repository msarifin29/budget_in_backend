package model

import "time"

type Account struct {
	UserId      string    `json:"user_id"`
	AccountId   string    `json:"account_id"`
	AccountName string    `json:"account_name"`
	Balance     float64   `json:"balance"`
	Cash        float64   `json:"cash"`
	Debts       float64   `json:"debts"`
	Savings     float64   `json:"savings"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type CreateAccountRequest struct {
	UserId      string  `json:"user_id" binding:"required"`
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance" binding:"required,min=2000"`
	Cash        float64 `json:"cash" binding:"required"`
	Currency    string  `json:"currency" binding:"required,currency"`
}
type GetAccountRequest struct {
	UserId    string `json:"user_id" binding:"required"`
	AccountId string `json:"account_id" binding:"required"`
}
type UpdateAccountBalance struct {
	UserId    string  `json:"user_id" binding:"required"`
	AccountId string  `json:"account_id" binding:"required"`
	Balance   float64 `json:"balance" binding:"required,min=2000"`
}
type UpdateAccountCash struct {
	UserId    string  `json:"user_id" binding:"required"`
	AccountId string  `json:"account_id" binding:"required"`
	Cash      float64 `json:"cash" binding:"required,min=2000"`
}
type UpdateAccountName struct {
	UserId      string `json:"user_id" binding:"required"`
	AccountId   string `json:"account_id" binding:"required"`
	AccountName string `json:"account_name" binding:"required"`
}

func NewAccount(account Account) *Account {
	return &Account{
		UserId:      account.UserId,
		AccountId:   account.AccountId,
		AccountName: account.AccountName,
		Balance:     account.Balance,
		Cash:        account.Cash,
		Debts:       account.Debts,
		Savings:     account.Savings,
		Currency:    account.Currency,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}
