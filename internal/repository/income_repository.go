package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type IncomeRepository interface {
	CreateIncome(ctx context.Context, tx *sql.Tx, income model.Income) (model.Income, error)
	GetIncomes(ctx context.Context, tx *sql.Tx, params model.GetIncomeParams) ([]model.Income, error)
	GetTotalIncomes(ctx context.Context, tx *sql.Tx, uid string, CategoryIncome string) (float64, error)
}

type IncomeRepositoryImpl struct{}

// GetTotalIncomes implements IncomeRepository.
func (*IncomeRepositoryImpl) GetTotalIncomes(ctx context.Context, tx *sql.Tx, uid string, CategoryIncome string) (float64, error) {
	var total float64
	script := `SELECT COUNT(*)from incomes where uid = ? && category_income = ?`
	err := tx.QueryRowContext(ctx, script, uid, CategoryIncome).Scan(&total)
	return total, err
}

// CreateIncome implements IncomeRepository.
func (*IncomeRepositoryImpl) CreateIncome(ctx context.Context, tx *sql.Tx, income model.Income) (model.Income, error) {
	script := `insert into incomes (uid,category_income,type_income,total) values (?,?,?,?)`
	result, errX := tx.ExecContext(ctx, script, income.Uid, income.CategoryIncome, income.TypeIncome, income.Total)
	if errX != nil {
		return model.Income{}, errX
	}
	lastId, err := result.LastInsertId()
	income.Id = float64(lastId)
	return income, err
}

// GetIncomes implements IncomeRepository.
func (*IncomeRepositoryImpl) GetIncomes(ctx context.Context, tx *sql.Tx, params model.GetIncomeParams) ([]model.Income, error) {
	script := `select * from incomes where uid = ? && category_income = ? order by id limit ? offset ?`
	rows, err := tx.QueryContext(ctx, script, params.Uid, params.CategoryIncome, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incomes := []model.Income{}
	var i model.Income
	update := zero.TimeFromPtr(&i.UpdatedAt)
	for rows.Next() {
		err := rows.Scan(
			&i.Uid,
			&i.Id,
			&i.CategoryIncome,
			&i.Total,
			&i.CreatedAt,
			&update,
			&i.TypeIncome,
		)
		if err != nil {
			return nil, err
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
