package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type ExpenseRepository interface {
	CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error)
	GetExpenseById(ctx context.Context, tx *sql.Tx, id float64) (model.Expense, error)
	UpdateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense, id float64) (model.Expense, error)
	DeleteExpense(ctx context.Context, tx *sql.Tx, id float64) error
	GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.ExpenseResponse, error)
	GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string, status string, expenseType string, Id string, CreatedAt string) (float64, error)
	GetExpenseThisMonth(ctx context.Context, tx *sql.Tx, uid string) (float64, error)
}

type ExpenseRepositoryImpl struct{}

// GetExpenseThisMonth implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenseThisMonth(ctx context.Context, tx *sql.Tx, uid string) (float64, error) {
	script := `
WITH expense_data AS (SELECT 
        TO_CHAR(created_at, 'YYYY-MM') AS month,
        uid,
        COALESCE(SUM(total), 0) AS total_expenses
    FROM  expenses 
    WHERE uid = $1 AND status = 'success' 
        AND EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE)
        AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE)
    GROUP BY month, uid)
SELECT 
    COALESCE(expense_data.month, TO_CHAR(CURRENT_DATE, 'YYYY-MM')) AS month,
    $1 AS uid,
    COALESCE(expense_data.total_expenses, 0) AS total_expenses
FROM (SELECT 1) AS dummy
LEFT JOIN expense_data ON TRUE
ORDER BY month ASC, uid ASC;`
	var total float64
	var userId, month string
	row := tx.QueryRowContext(ctx, script, uid)
	err := row.Scan(&month, &userId, &total)
	return total, err
}

// GetTotalExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetTotalExpenses(ctx context.Context, tx *sql.Tx, uid string, status string, expenseType string, Id string, CreatedAt string) (float64, error) {
	var total float64
	query := `SELECT COUNT(*) AS total
	FROM expenses e
	LEFT JOIN t_category_expenses t ON e.id = t.category_id
	WHERE uid = ? AND e.status = ?`
	if expenseType != "" {
		query += " AND e.expense_type LIKE ?"
	}
	if CreatedAt != "" {
		query += " AND DATE(e.created_at) = ?"
	}
	if Id != "" {
		query += " AND t.id = ?"
	}
	row := tx.QueryRowContext(ctx, query, append([]interface{}{uid, status}, optionalParams(expenseType, CreatedAt, Id)...)...)
	err := row.Scan(&total)
	return total, err
}

// GetExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.ExpenseResponse, error) {
	var err error
	query := `SELECT e.id, e.expense_type, e.total, e.notes, e.created_at, e.uid, e.status, e.transaction_id, t.category_id, t.id as t_id, t.title
	FROM expenses e
	LEFT JOIN t_category_expenses t ON e.id = t.category_id
	WHERE uid = ? AND status = ?`
	if params.ExpenseType != "" {
		query += " AND e.expense_type LIKE ?"
	}
	if params.CreatedAt != "" {
		query += " AND DATE(e.created_at) = ?"
	}
	if params.Id != "" {
		query += " AND t.id = ?"
	}
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	var rows *sql.Rows
	switch {
	case params.ExpenseType != "":
		rows, err = tx.QueryContext(ctx, query, params.Uid, params.Status, params.ExpenseType, params.Limit, params.Offset)
	case params.CreatedAt != "":
		rows, err = tx.QueryContext(ctx, query, params.Uid, params.Status, params.CreatedAt, params.Limit, params.Offset)
	case params.Id != "":
		rows, err = tx.QueryContext(ctx, query, params.Uid, params.Status, params.Id, params.Limit, params.Offset)
	default:
		rows, err = tx.QueryContext(ctx, query, params.Uid, params.Status, params.Limit, params.Offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.ExpenseResponse{}
	for rows.Next() {
		var i model.ExpenseResponse
		if err := rows.Scan(
			&i.Id, &i.ExpenseType, &i.Total, &i.Notes,
			&i.CreatedAt, &i.Uid, &i.Status,
			&i.TransactionId, &i.TCategory.CategoryId,
			&i.TCategory.Id, &i.TCategory.Title,
		); err != nil {
			return nil, err
		}
		if i.Notes == "" {
			i.Notes = ""
		}
		if i.TransactionId == "" {
			i.TransactionId = ""
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
	script := `select id, expense_type, total, notes, created_at, uid, status, transaction_id 
	from expenses where id = ? limit 1`

	rows := tx.QueryRowContext(ctx, script, id)

	var i model.Expense
	err := rows.Scan(
		&i.Id, &i.ExpenseType, &i.Total, &i.Notes,
		&i.CreatedAt, &i.Uid, &i.Category,
		&i.Status, &i.TransactionId,
	)
	if i.Notes == "" {
		i.Notes = ""
	}
	if i.TransactionId == "" {
		i.TransactionId = ""
	}
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

func optionalParams(params ...string) []interface{} {
	var values []interface{}
	for _, param := range params {
		if param != "" {
			values = append(values, param)
		}
	}
	return values
}
