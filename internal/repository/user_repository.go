package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error)
	GetUser(ctx context.Context, tx *sql.Tx, email string) (model.User, error)
	GetById(ctx context.Context, tx *sql.Tx, uid string) (model.User, error)
	UpdateUserName(ctx context.Context, tx *sql.Tx, user model.UpdateUserRequest) error
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error) {
	sqlScript := `INSERT INTO users (uid, username, email, password, type_user, balance, savings, cash, debts, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := tx.ExecContext(ctx, sqlScript,
		&user.Uid,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.TypeUser,
		&user.Balance,
		&user.Savings,
		&user.Cash,
		&user.Debts,
		&user.Currency)

	return model.User{
		Uid:      user.Uid,
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
		TypeUser: user.TypeUser,
		Balance:  user.Balance,
		Savings:  user.Savings,
		Cash:     user.Cash,
		Debts:    user.Debts,
		Currency: user.Currency,
	}, err
}

func (u *UserRepositoryImpl) GetUser(ctx context.Context, tx *sql.Tx, email string) (model.User, error) {
	sqlScript := `select uid, username, email, password, photo, created_at, updated_at, type_user, balance, savings, cash, debts, currency from users where email = ? limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, email)

	var i model.User
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := row.Scan(
		&i.Uid,
		&i.UserName,
		&i.Email,
		&i.Password,
		&i.Photo,
		&i.CreatedAt,
		&update,
		&i.TypeUser,
		&i.Balance,
		&i.Savings,
		&i.Cash,
		&i.Debts,
		&i.Currency,
	)
	return i, err
}

func (u *UserRepositoryImpl) GetById(ctx context.Context, tx *sql.Tx, uid string) (model.User, error) {
	sqlScript := `select uid, username, email, password, photo, created_at, updated_at, type_user, balance, savings, cash, debts, currency from users where uid = ? limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, uid)

	var i model.User
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := row.Scan(
		&i.Uid,
		&i.UserName,
		&i.Email,
		&i.Password,
		&i.Photo,
		&i.CreatedAt,
		&update,
		&i.TypeUser,
		&i.Balance,
		&i.Savings,
		&i.Cash,
		&i.Debts,
		&i.Currency,
	)
	return i, err
}

func (u *UserRepositoryImpl) UpdateUserName(ctx context.Context, tx *sql.Tx, user model.UpdateUserRequest) error {
	sqlScript := `update users set username = ? where uid = ?`
	_, err := tx.ExecContext(ctx, sqlScript, user.UserName, user.Uid)
	return err
}