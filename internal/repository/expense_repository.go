package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error)
	GetExpenseById(ctx context.Context, tx *sql.Tx, id float64) (model.Expense, error)
	UpdateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense, id float64) (model.Expense, error)
	DeleteExpense(ctx context.Context, tx *sql.Tx, id float64) error
	GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.ExpenseResponse, error)
	GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string, status string, expenseType string, Id int32) (float64, error)
	GetExpensesByMonth(ctx context.Context, tx *sql.Tx, params model.MonthlyParams) ([]model.Expense, error)
}

type ExpenseRepositoryImpl struct{}

// GetExpensesByMonth implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpensesByMonth(ctx context.Context, tx *sql.Tx, params model.MonthlyParams) ([]model.Expense, error) {
	script := `SELECT id,expense_type,total,created_at,uid,category,status
	from expenses where uid = ?
	AND YEAR(created_at) = ? AND MONTH(created_at) = ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, params.Year, params.Month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.Expense{}
	for rows.Next() {
		var i model.Expense
		if err := rows.Scan(
			&i.Id,
			&i.ExpenseType,
			&i.Total,
			&i.CreatedAt,
			&i.Uid,
			&i.Category,
			&i.Status,
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

// GetTotalExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string, status string, expenseType string, Id int32) (float64, error) {
	var total float64
	cId := categoryId(Id)
	script := `SELECT COUNT(*) AS total
	FROM expenses e
	LEFT JOIN t_category_expenses t ON e.id = t.category_id
	where uid = ? and e.status = ? and e.expense_type LIKE ? and t.id LIKE ?`
	err := tx.QueryRowContext(ctx, script, uid, status, "%"+expenseType+"%", "%"+cId+"%").Scan(&total)
	return total, err
}

// GetExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.ExpenseResponse, error) {
	cId := categoryId(params.Id)
	script := `select e.*, t.category_id, t.id as t_id, t.title
	from expenses e
	LEFT JOIN t_category_expenses t ON e.id = t.category_id
	where uid = ? and status = ? 
	and e.expense_type LIKE ? and t.id LIKE ?
	order by id limit ? offset ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, params.Status,
		"%"+params.ExpenseType+"%", "%"+cId+"%", params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.ExpenseResponse{}
	for rows.Next() {
		var i model.ExpenseResponse
		update := zero.TimeFromPtr(i.UpdatedAt)
		transaction := zero.StringFromPtr(&i.TransactionId)
		if err := rows.Scan(
			&i.Id, &i.ExpenseType,
			&i.Total, &i.Notes,
			&i.CreatedAt, &update,
			&i.Uid, &i.Category,
			&i.Status, &transaction,
			&i.TCategory.CategoryId,
			&i.TCategory.Id, &i.TCategory.Title,
		); err != nil {
			return nil, err
		}
		if i.Notes == "" {
			i.Notes = ""
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
	transaction := zero.StringFromPtr(&i.TransactionId)
	update := zero.TimeFromPtr(i.UpdatedAt)
	err := rows.Scan(
		&i.Id,
		&i.ExpenseType,
		&i.Total,
		&notes,
		&i.CreatedAt,
		&update,
		&i.Uid,
		&i.Category,
		&i.Status,
		&transaction,
	)
	return i, err
}

// CreateExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error) {
	script := `insert into expenses (expense_type,total,notes,uid,category,status,created_at,transaction_id) values (?,?,?,?,?,?,?,?)`
	result, errX := tx.ExecContext(ctx, script,
		&expense.ExpenseType, &expense.Total,
		&expense.Notes, &expense.Uid,
		&expense.Category, &expense.Status,
		&expense.CreatedAt, &expense.TransactionId)
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
	script := `update expenses set status = ? where id = ?`
	_, err := tx.ExecContext(ctx, script, expense.Status, id)
	return expense, err
}

func NewExpenseRepository() ExpenseRepository {
	return &ExpenseRepositoryImpl{}
}

func categoryId(id int32) string {
	if id == 0 {
		return ""
	}
	return fmt.Sprintf("%v", int(id))
}
