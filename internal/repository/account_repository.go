package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, tx *sql.Tx, account model.Account) (model.Account, error)
	GetAccountByUserId(ctx context.Context, tx *sql.Tx, account model.GetAccountRequest) (model.Account, error)
	GetAllAccount(ctx context.Context, tx *sql.Tx, userId string) ([]model.Account, error)
	UpdateAccountName(ctx context.Context, tx *sql.Tx, account model.UpdateAccountName) error
	UpdateAccountBalance(ctx context.Context, tx *sql.Tx, account model.UpdateAccountBalance) error
	UpdateAccountCash(ctx context.Context, tx *sql.Tx, account model.UpdateAccountCash) error
}

type AccountRepositoryImpl struct{}

// CreateAccount implements AccountRepository.
func (AccountRepositoryImpl) CreateAccount(ctx context.Context, tx *sql.Tx, account model.Account) (model.Account, error) {
	script := `insert into accounts (user_id,account_id,account_name,balance,cash,debts,savings,currency) values (?,?,?,?,?,?,?,?)`
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

// GetAccountByUserId implements AccountRepository.
func (AccountRepositoryImpl) GetAccountByUserId(ctx context.Context, tx *sql.Tx, account model.GetAccountRequest) (model.Account, error) {
	script := `select * from accounts where user_id = ? and account_id = ?`
	row := tx.QueryRowContext(ctx, script, account.UserId, account.AccountId)
	var i model.Account
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := row.Scan(
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
	return i, err
}

// GetAllAccount implements AccountRepository.
func (AccountRepositoryImpl) GetAllAccount(ctx context.Context, tx *sql.Tx, userId string) ([]model.Account, error) {
	script := `select * from accounts where user_id = ?`
	rows, err := tx.QueryContext(ctx, script, userId)
	if err != nil {
		return nil, err
	}
	var accounts []model.Account
	for rows.Next() {
		var i model.Account
		update := zero.TimeFromPtr(&i.UpdatedAt)
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
	script := `update accounts set balance = ? where user_id = ? and account_id = ?`
	_, err := tx.ExecContext(ctx, script, account.Balance, account.UserId, account.AccountId)
	return err
}

// UpdateAccountCash implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountCash(ctx context.Context, tx *sql.Tx, account model.UpdateAccountCash) error {
	script := `update accounts set cash = ? where user_id = ? and account_id = ?`
	_, err := tx.ExecContext(ctx, script, account.Cash, account.UserId, account.AccountId)
	return err
}

// UpdateAccountName implements AccountRepository.
func (AccountRepositoryImpl) UpdateAccountName(ctx context.Context, tx *sql.Tx, account model.UpdateAccountName) error {
	script := `update accounts set account_name = ? where user_id = ? and account_id = ?`
	_, err := tx.ExecContext(ctx, script, account.AccountName, account.UserId, account.AccountId)
	return err
}

func NewAccountRepository() AccountRepository {
	return AccountRepositoryImpl{}
}
