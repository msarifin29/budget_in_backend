package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, tx *sql.Tx, account model.Account) (model.Account, error)
	GetAccountByAccountId(ctx context.Context, tx *sql.Tx, account model.GetAccountRequest) (model.Account, error)
	GetAllAccount(ctx context.Context, tx *sql.Tx, userId string) ([]model.Account, error)
	UpdateAccountName(ctx context.Context, tx *sql.Tx, account model.UpdateAccountName) error
	UpdateAccountBalance(ctx context.Context, tx *sql.Tx, account model.UpdateAccountBalance) error
	UpdateAccountCash(ctx context.Context, tx *sql.Tx, account model.UpdateAccountCash) error
	UpdateAccountDebts(ctx context.Context, tx *sql.Tx, account model.UpdateAccountDebts) error
	UpdateMaxBudget(ctx context.Context, tx *sql.Tx, account model.UpdateMaxBudgetRequest) error
	GetAccountByUserId(ctx context.Context, tx *sql.Tx, userId string) (model.Account, error)
}

type AccountRepositoryImpl struct{}

// GetAccountByUserId implements AccountRepository.
func (AccountRepositoryImpl) GetAccountByUserId(ctx context.Context, tx *sql.Tx, userId string) (model.Account, error) {
	script := `select user_id, account_id, account_name, balance, cash, debts, savings, currency,max_budget, created_at, updated_at
	 from accounts where user_id = $1`
	row := tx.QueryRowContext(ctx, script, userId)
	var i model.Account
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(
		&i.UserId, &i.AccountId,
		&i.AccountName, &i.Balance,
		&i.Cash, &i.Debts,
		&i.Savings, &i.Currency,
		&i.MaxBudget, &i.CreatedAt,
		&update,
	)
	return i, err
}

// UpdateMaxBudget implements AccountRepository.
func (AccountRepositoryImpl) UpdateMaxBudget(ctx context.Context, tx *sql.Tx, account model.UpdateMaxBudgetRequest) error {
	script := `update accounts set max_budget = $1 where account_id = $2 and user_id = $3`
	_, err := tx.ExecContext(ctx, script, account.MaxBudget, account.AccountId, account.Uid)
	return err
}

// UpdateAccountDebts implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountDebts(ctx context.Context, tx *sql.Tx, account model.UpdateAccountDebts) error {
	script := `update accounts set debts = $1 where account_id = $2`
	_, err := tx.ExecContext(ctx, script, account.Debts, account.AccountId)
	return err
}

// CreateAccount implements AccountRepository.
func (AccountRepositoryImpl) CreateAccount(ctx context.Context, tx *sql.Tx, account model.Account) (model.Account, error) {
	script := `insert into accounts (user_id,account_id,account_name,balance,cash,debts,savings,currency) values ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := tx.ExecContext(ctx, script,
		account.UserId,
		account.AccountId,
		account.AccountName,
		account.Balance,
		account.Cash,
		account.Debts,
		account.Savings,
		account.Currency,
	)
	return account, err
}

// GetAccountByAccountId implements AccountRepository.
func (AccountRepositoryImpl) GetAccountByAccountId(ctx context.Context, tx *sql.Tx, account model.GetAccountRequest) (model.Account, error) {
	script := `select user_id, account_id, account_name, balance, cash, debts, savings, currency,max_budget, created_at, updated_at
	 from accounts where account_id = $1`
	row := tx.QueryRowContext(ctx, script, account.AccountId)
	var i model.Account
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(
		&i.UserId, &i.AccountId,
		&i.AccountName, &i.Balance,
		&i.Cash, &i.Debts,
		&i.Savings, &i.Currency,
		&i.MaxBudget, &i.CreatedAt,
		&update,
	)
	return i, err
}

// GetAllAccount implements AccountRepository.
func (AccountRepositoryImpl) GetAllAccount(ctx context.Context, tx *sql.Tx, userId string) ([]model.Account, error) {
	script := `select user_id, account_id, account_name, balance, cash, debts, savings, currency, created_at, updated_at 
	from accounts where user_id = $1`
	rows, err := tx.QueryContext(ctx, script, userId)
	if err != nil {
		return nil, err
	}
	var accounts []model.Account
	for rows.Next() {
		var i model.Account
		update := zero.TimeFromPtr(i.UpdatedAt)
		err := rows.Scan(
			&i.UserId,
			&i.AccountId,
			&i.AccountName,
			&i.Balance,
			&i.Cash,
			&i.Debts,
			&i.Savings,
			&i.Currency,
			&i.CreatedAt,
			&update,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

// UpdateAccountBalance implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountBalance(ctx context.Context, tx *sql.Tx, account model.UpdateAccountBalance) error {
	script := `update accounts set balance = $1 where account_id = $2`
	_, err := tx.ExecContext(ctx, script, account.Balance, account.AccountId)
	return err
}

// UpdateAccountCash implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountCash(ctx context.Context, tx *sql.Tx, account model.UpdateAccountCash) error {
	script := `update accounts set cash = $1 where account_id = $2`
	_, err := tx.ExecContext(ctx, script, account.Cash, account.AccountId)
	return err
}

// UpdateAccountName implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountName(ctx context.Context, tx *sql.Tx, account model.UpdateAccountName) error {
	script := `update accounts set account_name = $1 where account_id = $2`
	_, err := tx.ExecContext(ctx, script, account.AccountName, account.AccountId)
	return err
}

func NewAccountRepository() AccountRepository {
	return AccountRepositoryImpl{}
}
