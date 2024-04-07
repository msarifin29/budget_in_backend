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
	GetUserAccount(ctx context.Context, tx *sql.Tx, uid string) (model.AccountUser, error)
	UpdateUserName(ctx context.Context, tx *sql.Tx, user model.UpdateUserRequest) error
	GetUserByEmail(ctx context.Context, tx *sql.Tx, req model.EmailUserRequest) (string, string, error)
	UpdatePassword(ctx context.Context, tx *sql.Tx, email string, newPassword string) (bool, error)
	NonActivatedUser(ctx context.Context, tx *sql.Tx, uid string, status string) (bool, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// NonActivatedUser implements UserRepository.
func (*UserRepositoryImpl) NonActivatedUser(ctx context.Context, tx *sql.Tx, uid string, status string) (bool, error) {
	sqlScript := `update users set status = ? where uid = ?`
	_, err := tx.ExecContext(ctx, sqlScript, status, uid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUserAccount implements UserRepository.
func (*UserRepositoryImpl) GetUserAccount(ctx context.Context, tx *sql.Tx, uid string) (model.AccountUser, error) {
	sqlScript := `SELECT accounts.user_id as uid,accounts.account_id, users.username,users.email, accounts.balance,accounts.cash,accounts.debts,
	accounts.savings,accounts.currency,accounts.created_at,accounts.updated_at
	FROM users CROSS JOIN accounts
	on users.uid = accounts.user_id 
	WHERE users.uid = ?`
	row := tx.QueryRowContext(ctx, sqlScript, uid)

	var i model.AccountUser
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := row.Scan(
		&i.Uid, &i.AccountId, &i.UserName,
		&i.Email, &i.Balance, &i.Cash,
		&i.Debts, &i.Savings, &i.Currency,
		&i.CreatedAt, &update,
	)
	return i, err
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error) {
	sqlScript := `INSERT INTO users (uid, username, email, password, type_user, balance, savings, cash, debts, currency) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := tx.ExecContext(ctx, sqlScript,
		&user.Uid, &user.UserName, &user.Email,
		&user.Password, &user.TypeUser, &user.Balance,
		&user.Savings, &user.Cash, &user.Debts,
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
	sqlScript := `select uid, username, email, password, photo, created_at, updated_at, type_user, balance, savings, cash, debts, currency 
	from users where email = ? and status = "active" limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, email)

	var i model.User
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(
		&i.Uid, &i.UserName,
		&i.Email, &i.Password,
		&i.Photo, &i.CreatedAt,
		&update, &i.TypeUser,
		&i.Balance, &i.Savings, &i.Cash,
		&i.Debts, &i.Currency,
	)
	return i, err
}

func (u *UserRepositoryImpl) GetById(ctx context.Context, tx *sql.Tx, uid string) (model.User, error) {
	sqlScript := `select * from users where uid = ? and status = "active" limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, uid)

	var i model.User
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(
		&i.Uid, &i.UserName, &i.Email,
		&i.Password, &i.Photo, &i.CreatedAt,
		&update, &i.TypeUser, &i.Balance,
		&i.Savings, &i.Cash, &i.Debts,
		&i.Currency, &i.Status,
	)
	return i, err
}

func (u *UserRepositoryImpl) UpdateUserName(ctx context.Context, tx *sql.Tx, user model.UpdateUserRequest) error {
	sqlScript := `update users set username = ? where uid = ? and status = "active"`
	_, err := tx.ExecContext(ctx, sqlScript, user.UserName, user.Uid)
	return err
}

// GetUserByEmail implements UserRepository.
func (*UserRepositoryImpl) GetUserByEmail(ctx context.Context, tx *sql.Tx, req model.EmailUserRequest) (string, string, error) {
	sqlScript := `select email, username from users where email = ? and status = "active" limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, req.Email)
	var email, username string
	err := row.Scan(&email, &username)
	return email, username, err
}

// UpdatePassword implements UserRepository.
func (*UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, email string, newPassword string) (bool, error) {
	sqlScript := `update users set password = ? where email = ? and status = "active"`
	_, err := tx.ExecContext(ctx, sqlScript, newPassword, email)
	if err != nil {
		return false, err
	}
	return true, nil
}
