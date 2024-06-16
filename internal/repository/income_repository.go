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
	where uid = $1 and type_income LIKE $2`
	param := []interface{}{uid, `%` + typeIncome + `%`}
	if len(param) == 3 {
		script += " and i.id = $3"
		param = append(param, cId)
	}
	err := tx.QueryRowContext(ctx, script, param...).Scan(&total)
	return total, err
}

// CreateIncome implements IncomeRepository.
func (*IncomeRepositoryImpl) CreateIncome(ctx context.Context, tx *sql.Tx, income model.Income) (model.Income, error) {
	script := `insert into incomes (uid,type_income,total,created_at,transaction_id,account_id,bank_name,bank_id) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	var id int64
	errX := tx.QueryRowContext(ctx, script, &income.Uid, &income.TypeIncome, &income.Total, &income.CreatedAt,
		&income.TransactionId, &income.AccountId, &income.BankName, &income.BankId).Scan(&id)
	if errX != nil {
		return model.Income{}, errX
	}
	income.Id = float64(id)
	return income, errX
}

// GetIncomes implements IncomeRepository.
func (*IncomeRepositoryImpl) GetIncomes(ctx context.Context, tx *sql.Tx, params model.GetIncomeParams) ([]model.IncomeResponse, error) {
	cId := categoryId(params.CategoryId)
	script := `
	select i.uid, i.id, i.total, i.created_at, i.updated_at, i.type_income, i.transaction_id, t.category_id, t.id as t_id, t.title, i.account_id,
	i.bank_name, i.bank_id from incomes i
	LEFT JOIN t_category_incomes t ON i.id = t.category_id
	where uid = $1 and i.type_income LIKE $2`

	param := []interface{}{params.Uid, `%` + params.TypeIncome + `%`}
	if cId != "" {
		script += " and t.id = $3"
		param = append(param, cId)
	}
	if len(param) == 3 {
		script += " ORDER BY id DESC LIMIT $4 OFFSET $5"
		param = append(param, params.Limit, params.Offset)
	} else {
		script += " ORDER BY id DESC LIMIT $3 OFFSET $4"
		param = append(param, params.Limit, params.Offset)
	}
	rows, err := tx.QueryContext(ctx, script, param...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incomes := []model.IncomeResponse{}
	var i model.IncomeResponse
	update := zero.TimeFromPtr(i.UpdatedAt)
	for rows.Next() {
		err := rows.Scan(
			&i.Uid, &i.Id, &i.Total, &i.CreatedAt, &update, &i.TypeIncome, &i.TransactionId,
			&i.TCategory.CategoryId, &i.TCategory.Id, &i.TCategory.Title, &i.AccountId, &i.BankName, &i.BankId,
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
