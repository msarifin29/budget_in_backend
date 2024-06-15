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
	ResetPassword(ctx context.Context, tx *sql.Tx, params model.ResetPasswordRequest) (bool, error)
	NonActivatedUser(ctx context.Context, tx *sql.Tx, uid string, email string) (bool, error)
	GetEmailUser(ctx context.Context, tx *sql.Tx) ([]model.User, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// GetEmailUser implements UserRepository.
func (*UserRepositoryImpl) GetEmailUser(ctx context.Context, tx *sql.Tx) ([]model.User, error) {
	sqlScript := `select email from users`
	rows, err := tx.QueryContext(ctx, sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []model.User{}
	for rows.Next() {
		var i model.User
		err := rows.Scan(&i.Email)
		if err != nil {
			return nil, err
		}
		if i.Email == "" {
			i.Email = ""
		}
		users = append(users, i)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return users, nil

}

// ResetPassword implements UserRepository.
func (*UserRepositoryImpl) ResetPassword(ctx context.Context, tx *sql.Tx, params model.ResetPasswordRequest) (bool, error) {
	sqlScript := `update users set password = $1 where uid = $2 and status = 'active'`
	_, err := tx.ExecContext(ctx, sqlScript, params.NewPassword, params.Uid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// NonActivatedUser implements UserRepository.
func (*UserRepositoryImpl) NonActivatedUser(ctx context.Context, tx *sql.Tx, uid string, email string) (bool, error) {
	sqlScript := `update users set email = $1 where uid = $2`
	_, err := tx.ExecContext(ctx, sqlScript, email, uid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUserAccount implements UserRepository.
func (*UserRepositoryImpl) GetUserAccount(ctx context.Context, tx *sql.Tx, uid string) (model.AccountUser, error) {
	sqlScript := `SELECT accounts.user_id as uid,accounts.account_id, users.username, accounts.account_name, users.email, accounts.balance,accounts.cash,accounts.debts,
	accounts.savings,accounts.currency,accounts.created_at,accounts.updated_at
	FROM users JOIN accounts
	on users.uid = accounts.user_id 
	WHERE users.uid = $1`
	row := tx.QueryRowContext(ctx, sqlScript, uid)

	var i model.AccountUser
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := row.Scan(
		&i.Uid, &i.AccountId, &i.UserName, &i.AccountName,
		&i.Email, &i.Balance, &i.Cash,
		&i.Debts, &i.Savings, &i.Currency,
		&i.CreatedAt, &update,
	)
	return i, err
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user model.User) (model.User, error) {
	sqlScript := `INSERT INTO users (uid, username, email, password, type_user) VALUES ($1, $2, $3, $4, $5);`
	_, err := tx.ExecContext(ctx, sqlScript, &user.Uid, &user.UserName, &user.Email, &user.Password, &user.TypeUser)

	return model.User{
		Uid:      user.Uid,
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
		TypeUser: user.TypeUser,
	}, err
}

func (u *UserRepositoryImpl) GetUser(ctx context.Context, tx *sql.Tx, email string) (model.User, error) {
	sqlScript := `select uid, username, email, password, photo, created_at, updated_at, type_user 
	from users where email = $1 and status = 'active' limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, email)

	var i model.User
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(&i.Uid, &i.UserName, &i.Email, &i.Password, &i.Photo, &i.CreatedAt, &update, &i.TypeUser)
	return i, err
}

func (u *UserRepositoryImpl) GetById(ctx context.Context, tx *sql.Tx, uid string) (model.User, error) {
	sqlScript := `select uid, username, email, password, photo, created_at, updated_at, type_user, status
	 from users where uid = $1 and status = 'active' limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, uid)

	var i model.User
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := row.Scan(&i.Uid, &i.UserName, &i.Email, &i.Password, &i.Photo, &i.CreatedAt, &update, &i.TypeUser, &i.Status)
	return i, err
}

func (u *UserRepositoryImpl) UpdateUserName(ctx context.Context, tx *sql.Tx, user model.UpdateUserRequest) error {
	sqlScript := `update users set username = $1 where uid = $2 and status = 'active'`
	_, err := tx.ExecContext(ctx, sqlScript, user.UserName, user.Uid)
	return err
}

// GetUserByEmail implements UserRepository.
func (*UserRepositoryImpl) GetUserByEmail(ctx context.Context, tx *sql.Tx, req model.EmailUserRequest) (string, string, error) {
	sqlScript := `select email, username from users where email = $1 and status = 'active' limit 1`
	row := tx.QueryRowContext(ctx, sqlScript, req.Email)
	var email, username string
	err := row.Scan(&email, &username)
	return email, username, err
}

// UpdatePassword implements UserRepository.
func (*UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *sql.Tx, email string, newPassword string) (bool, error) {
	sqlScript := `update users set password = $1 where email = $2 and status = 'active'`
	_, err := tx.ExecContext(ctx, sqlScript, newPassword, email)
	if err != nil {
		return false, err
	}
	return true, nil
}
