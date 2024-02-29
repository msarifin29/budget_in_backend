package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error)
	GetExpenseById(ctx context.Context, tx *sql.Tx, id float64) (model.Expense, error)
	UpdateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense, id float64) (model.Expense, error)
	DeleteExpense(ctx context.Context, tx *sql.Tx, id float64) error
	GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.Expense, error)
	GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
}

type ExpenseRepositoryImpl struct{}

// GetTotalExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
	var total float64
	script := `SELECT COUNT(*)from expenses where uid = ?`
	err := tx.QueryRowContext(ctx, script, uid).Scan(&total)
	return total, err
}

// GetExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.Expense, error) {
	script := `select id, expense_type, total, notes, created_at, updated_at, uid from expenses where uid = ? order by id limit ? offset ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.Expense{}
	for rows.Next() {
		var i model.Expense
		notes := zero.StringFromPtr(&i.Notes)
		update := zero.TimeFromPtr(&i.UpdatedAt)
		if err := rows.Scan(
			&i.Id,
			&i.ExpenseType,
			&i.Total,
			&notes,
			&i.CreatedAt,
			&update,
			&i.Uid,
		); err != nil {
			return nil, err
		}
		expenses = append(expenses, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

// GetExpenseById implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenseById(ctx context.Context, tx *sql.Tx, id float64) (model.Expense, error) {
	script := `select * from expenses where id = ? limit 1`

	rows := tx.QueryRowContext(ctx, script, id)

	var i model.Expense
	notes := zero.StringFromPtr(&i.Notes)
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := rows.Scan(
		&i.Id,
		&i.ExpenseType,
		&i.Total,
		&notes,
		&i.CreatedAt,
		&update,
		&i.Uid,
	)
	return i, err
}

// CreateExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error) {
	script := `insert into expenses (expense_type,total,notes,uid) values (?,?,?,?)`
	result, errX := tx.ExecContext(ctx, script, &expense.ExpenseType, &expense.Total, &expense.Notes, &expense.Uid)
	if errX != nil {
		return model.Expense{}, errX
	}
	id, err := result.LastInsertId()
	expense.Id = float64(id)
	return expense, err
}

// DeleteExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) DeleteExpense(ctx context.Context, tx *sql.Tx, id float64) error {
	script := `delete from expenses where id = ?`
	_, err := tx.ExecContext(ctx, script, id)
	return err
}

// UpdateExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) UpdateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense, id float64) (model.Expense, error) {
	script := `update expenses set expense_type = ?, total = ?, notes = ? where id = ?`
	_, err := tx.ExecContext(ctx, script, expense.ExpenseType, expense.Total, expense.Notes, id)
	return expense, err
}

func NewExpenseRepository() ExpenseRepository {
	return &ExpenseRepositoryImpl{}
}
