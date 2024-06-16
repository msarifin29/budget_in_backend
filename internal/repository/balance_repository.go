package repository

// import (
// 	"context"
// 	"database/sql"
// )

// type BalanceRepository interface {
// 	GetBalance(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
// 	SetBalance(ctx context.Context, tx *sql.Tx, uid string, i float64) error
// 	GetSaving(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
// 	SetSaving(ctx context.Context, tx *sql.Tx, uid string, i float64) error
// 	GetCash(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
// 	SetCash(ctx context.Context, tx *sql.Tx, uid string, i float64) error
// 	GetDebt(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
// 	SetDebt(ctx context.Context, tx *sql.Tx, uid string, i float64) error
// }
// type BalanceRepositoryImpl struct {
// }

// // GetBalance implements BalanceRepository.
// func (*BalanceRepositoryImpl) GetBalance(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
// 	var n float64
// 	row := tx.QueryRowContext(ctx, `select balance from users where uid = $1`, uid)
// 	err := row.Scan(&n)
// 	return n, err
// }

// // GetCash implements BalanceRepository.
// func (*BalanceRepositoryImpl) GetCash(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
// 	var n float64
// 	row := tx.QueryRowContext(ctx, `select cash from users where uid = $1`, uid)
// 	err := row.Scan(&n)
// 	return n, err
// }

// // GetDebt implements BalanceRepository.
// func (*BalanceRepositoryImpl) GetDebt(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
// 	var n float64
// 	row := tx.QueryRowContext(ctx, `select debts from users where uid = $1`, uid)
// 	err := row.Scan(&n)
// 	return n, err
// }

// // GetSaving implements BalanceRepository.
// func (*BalanceRepositoryImpl) GetSaving(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
// 	var n float64
// 	row := tx.QueryRowContext(ctx, `select savings from users where uid = $1`, uid)
// 	err := row.Scan(&n)
// 	return n, err
// }

// // SetBalance implements BalanceRepository.
// func (*BalanceRepositoryImpl) SetBalance(ctx context.Context, tx *sql.Tx, uid string, i float64) error {
// 	_, err := tx.ExecContext(ctx, `update users set balance = ? where uid = $1`, i, uid)
// 	return err
// }

// // SetCash implements BalanceRepository.
// func (*BalanceRepositoryImpl) SetCash(ctx context.Context, tx *sql.Tx, uid string, i float64) error {
// 	_, err := tx.ExecContext(ctx, `update users set cash = ? where uid = $1`, i, uid)
// 	return err
// }

// // SetDebt implements BalanceRepository.
// func (*BalanceRepositoryImpl) SetDebt(ctx context.Context, tx *sql.Tx, uid string, i float64) error {
// 	_, err := tx.ExecContext(ctx, `update users set debts = ? where uid = $1`, i, uid)
// 	return err

// }

// // SetSaving implements BalanceRepository.
// func (*BalanceRepositoryImpl) SetSaving(ctx context.Context, tx *sql.Tx, uid string, i float64) error {
// 	_, err := tx.ExecContext(ctx, `update users set savings = ? where uid = $1`, i, uid)
// 	return err
// }

// func NewBalanceRepository() BalanceRepository {
// 	return &BalanceRepositoryImpl{}
// }
