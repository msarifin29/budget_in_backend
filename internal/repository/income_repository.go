package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type IncomeRepository interface {
	CreateIncome(ctx context.Context, tx *sql.Tx, income model.Income) (model.Income, error)
	GetIncomes(ctx context.Context, tx *sql.Tx, params model.GetIncomeParams) ([]model.IncomeResponse, error)
	GetTotalIncomes(ctx context.Context, tx *sql.Tx, uid string, Id int32, typeIncome string) (float64, error)
	GetIncomesByMonth(ctx context.Context, tx *sql.Tx, params model.MonthlyParams) ([]model.Income, error)
}

type IncomeRepositoryImpl struct{}

// Deprecated: GetIncomesByMonth implements IncomeRepository.
func (*IncomeRepositoryImpl) GetIncomesByMonth(ctx context.Context, tx *sql.Tx, params model.MonthlyParams) ([]model.Income, error) {
	script := `SELECT uid,id,category_income,total,created_at,type_income
	from incomes where uid = ?
	AND YEAR(created_at) = ? AND MONTH(created_at) = ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, params.Year, params.Month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incomes := []model.Income{}
	for rows.Next() {
		var i model.Income
		if err := rows.Scan(
			&i.Uid,
			&i.Id,
			&i.CategoryIncome,
			&i.Total,
			&i.CreatedAt,
			&i.TypeIncome,
		); err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}

// GetTotalIncomes implements IncomeRepository.
func (*IncomeRepositoryImpl) GetTotalIncomes(ctx context.Context, tx *sql.Tx, uid string, Id int32, typeIncome string) (float64, error) {
	var total float64
	cId := categoryId(Id)
	script := `SELECT COUNT(*)as total from incomes i
	LEFT JOIN t_category_incomes t ON i.id = t.category_id
	where uid = ? && t.id LIKE ? and type_income LIKE ?`
	err := tx.QueryRowContext(ctx, script, uid, `%`+cId+`%`, `%`+typeIncome+`%`).Scan(&total)
	return total, err
}

// CreateIncome implements IncomeRepository.
func (*IncomeRepositoryImpl) CreateIncome(ctx context.Context, tx *sql.Tx, income model.Income) (model.Income, error) {
	script := `insert into incomes (uid,category_income,type_income,total,created_at,transaction_id) values (?,?,?,?,?,?)`
	result, errX := tx.ExecContext(ctx, script, income.Uid, income.CategoryIncome, income.TypeIncome, income.Total, income.CreatedAt, income.TransactionId)
	if errX != nil {
		return model.Income{}, errX
	}
	lastId, err := result.LastInsertId()
	income.Id = float64(lastId)
	return income, err
}

// GetIncomes implements IncomeRepository.
func (*IncomeRepositoryImpl) GetIncomes(ctx context.Context, tx *sql.Tx, params model.GetIncomeParams) ([]model.IncomeResponse, error) {
	cId := categoryId(params.CategoryId)
	script := `
	select i.*, t.category_id, t.id as t_id, t.title
	from incomes i
	LEFT JOIN t_category_incomes t ON i.id = t.category_id
	where uid = ? and i.type_income LIKE ? and t.id LIKE ?
	order by id desc limit ? offset ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, `%`+params.TypeIncome+`%`, `%`+cId+`%`, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incomes := []model.IncomeResponse{}
	var i model.IncomeResponse
	update := zero.TimeFromPtr(i.UpdatedAt)
	for rows.Next() {
		err := rows.Scan(
			&i.Uid, &i.Id,
			&i.CategoryIncome, &i.Total,
			&i.CreatedAt, &update,
			&i.TypeIncome, &i.TransactionId,
			&i.TCategory.CategoryId,
			&i.TCategory.Id, &i.TCategory.Title,
		)
		if err != nil {
			return nil, err
		}
		if i.TransactionId == "" {
			i.TransactionId = ""
		}
		incomes = append(incomes, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}

func NewIncomeRepository() IncomeRepository {
	return &IncomeRepositoryImpl{}
}
