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
	WHERE uid = $1 AND e.status = $2`
	param := []interface{}{uid, status}
	if expenseType != "" {
		query += " AND e.expense_type LIKE $3"
		param = append(param, "%"+expenseType+"%")
	}
	if CreatedAt != "" {
		if len(param) == 3 {
			query += " AND DATE(e.created_at) = $4"
			param = append(param, CreatedAt)
		} else {
			query += " AND DATE(e.created_at) = $3"
			param = append(param, CreatedAt)
		}
	}
	if Id != "" {
		if len(param) == 4 {
			query += " AND t.id = $5"
			param = append(param, Id)
		} else if len(param) == 3 {
			query += " AND t.id = $4"
			param = append(param, Id)
		} else {
			query += " AND t.id = $3"
			param = append(param, Id)
		}
	}
	row := tx.QueryRowContext(ctx, query, param...)
	err := row.Scan(&total)
	return total, err
}

// GetExpenses implements ExpenseRepository.
func (*ExpenseRepositoryImpl) GetExpenses(ctx context.Context, tx *sql.Tx, params model.GetExpenseParams) ([]model.ExpenseResponse, error) {
	var err error
	query := `SELECT e.id, e.expense_type, e.total, e.notes, e.created_at, e.uid, e.status, e.transaction_id, 
	t.category_id, t.id as t_id, t.title, e.account_id, e.bank_name, e.bank_id
	FROM expenses e
	LEFT JOIN t_category_expenses t ON e.id = t.category_id
	WHERE uid = $1 AND e.status = $2`
	param := []interface{}{params.Uid, params.Status}
	if params.ExpenseType != "" {
		query += " AND e.expense_type LIKE $3"
		param = append(param, "%"+params.ExpenseType+"%")
	}
	if params.CreatedAt != "" {
		if len(param) == 3 {
			query += " AND DATE(e.created_at) = $4"
			param = append(param, params.CreatedAt)
		} else {
			query += " AND DATE(e.created_at) = $3"
			param = append(param, params.CreatedAt)
		}
	}
	if params.Id != "" {
		if len(param) == 4 {
			query += " AND t.id = $5"
			param = append(param, params.Id)
		} else if len(param) == 3 {
			query += " AND t.id = $4"
			param = append(param, params.Id)
		} else {
			query += " AND t.id = $3"
			param = append(param, params.Id)
		}
	}
	if len(param) == 3 {
		query += " ORDER BY id DESC LIMIT $4 OFFSET $5"
		param = append(param, params.Limit, params.Offset)
	} else {
		query += " ORDER BY id DESC LIMIT $3 OFFSET $4"
		param = append(param, params.Limit, params.Offset)
	}
	rows, err := tx.QueryContext(ctx, query, param...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.ExpenseResponse{}
	for rows.Next() {
		var i model.ExpenseResponse
		if err := rows.Scan(
			&i.Id, &i.ExpenseType, &i.Total, &i.Notes, &i.CreatedAt, &i.Uid, &i.Status,
			&i.TransactionId, &i.TCategory.CategoryId,
			&i.TCategory.Id, &i.TCategory.Title, &i.AccountId, &i.BankName, &i.BankId,
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
	script := `select id, expense_type, total, notes, created_at, uid, status, transaction_id, account_id, bank_name, bank_id 
	from expenses where id = $1 limit 1`

	rows := tx.QueryRowContext(ctx, script, id)

	var i model.Expense
	err := rows.Scan(&i.Id, &i.ExpenseType, &i.Total, &i.Notes, &i.CreatedAt, &i.Uid, &i.Status, &i.TransactionId, &i.AccountId, &i.BankName, &i.BankId)
	if i.Notes == "" {
		i.Notes = ""
	}
	return i, err
}

// CreateExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) CreateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense) (model.Expense, error) {
	var id int
	script := `insert into expenses (expense_type,total,notes,uid,status,created_at,transaction_id, account_id, bank_name, bank_id,c_id) 
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`
	errX := tx.QueryRowContext(ctx, script, &expense.ExpenseType, &expense.Total, &expense.Notes, &expense.Uid,
		&expense.Status, &expense.CreatedAt, &expense.TransactionId, &expense.AccountId, &expense.BankName, &expense.BankId, expense.Cid).Scan(&id)
	if errX != nil {
		return model.Expense{}, errX
	}
	expense.Id = float64(id)
	return expense, errX
}

// DeleteExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) DeleteExpense(ctx context.Context, tx *sql.Tx, id float64) error {
	script := `delete from expenses where id = ?`
	_, err := tx.ExecContext(ctx, script, id)
	return err
}

// UpdateExpense implements ExpenseRepository.
func (*ExpenseRepositoryImpl) UpdateExpense(ctx context.Context, tx *sql.Tx, expense model.Expense, id float64) (model.Expense, error) {
	script := `update expenses set status = $1 where id = $2`
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
