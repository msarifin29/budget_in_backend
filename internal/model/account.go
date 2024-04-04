package model

import "time"

type Account struct {
	UserId      string     `json:"user_id"`
	AccountId   string     `json:"account_id"`
	AccountName string     `json:"account_name"`
	Balance     float64    `json:"balance"`
	Cash        float64    `json:"cash"`
	Debts       float64    `json:"debts"`
	Savings     float64    `json:"savings"`
	Currency    string     `json:"currency"`
	MaxBudget   float64    `json:"max_budget"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
type CreateAccountRequest struct {
	UserId      string  `json:"user_id" binding:"required"`
	AccountName string  `json:"account_name"`
	Balance     float64 `json:"balance" binding:"required,min=2000"`
	Cash        float64 `json:"cash" binding:"required"`
	Currency    string  `json:"currency" binding:"required,currency"`
}
type GetAccountRequest struct {
	AccountId string `json:"account_id" binding:"required"`
}
type UpdateAccountBalance struct {
	AccountId string  `json:"account_id" binding:"required"`
	Balance   float64 `json:"balance" binding:"required,min=2000"`
}
type UpdateAccountCash struct {
	AccountId string  `json:"account_id" binding:"required"`
	Cash      float64 `json:"cash" binding:"required,min=2000"`
}
type UpdateAccountDebts struct {
	AccountId string  `json:"account_id" binding:"required"`
	Debts     float64 `json:"debts" binding:"required,min=2000"`
}
type UpdateAccountName struct {
	AccountId   string `json:"account_id" binding:"required"`
	AccountName string `json:"account_name" binding:"required"`
}

type UpdateMaxBudgetRequest struct {
	Uid       string  `json:"uid" binding:"required"`
	AccountId string  `json:"account_id" binding:"required"`
	MaxBudget float64 `json:"max_budget" binding:"min=10000"`
}
type GetMaxBudgetParam struct {
	Uid       string `json:"uid"`
	AccountId string `json:"account_id"`
}
type GetMaxBudgetRequest struct {
	Uid       string `form:"uid" binding:"required"`
	AccountId string `form:"account_id" binding:"required"`
}
type MaxBudgetResponse struct {
	Uid          string  `json:"uid"`
	AccountId    string  `json:"account_id"`
	MaxBudget    float64 `json:"max_budget"`
	TotalExpense float64 `json:"total_expense"`
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
		MaxBudget:   account.MaxBudget,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}
