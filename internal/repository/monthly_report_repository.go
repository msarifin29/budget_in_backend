package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type MonthlyReportRepository interface {
	MonthlyReportExpensesDetail(ctx context.Context, tx *sql.Tx, param model.ParamMonthlyReportDetail) ([]model.MonthlyXDetail, error)
	MonthlyReportIncomesDetail(ctx context.Context, tx *sql.Tx, param model.ParamMonthlyReportDetail) ([]model.MonthlyIDetail, error)
	GetMonthlyIncomeReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
	GetMonthlyExpenseReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
}

type MonthlyReportRepositoryImpl struct{}

// MonthlyReportIncomesDetail implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) MonthlyReportIncomesDetail(ctx context.Context, tx *sql.Tx, param model.ParamMonthlyReportDetail) ([]model.MonthlyIDetail, error) {
	script := `
SELECT 
    DATE_FORMAT(i.created_at, '%Y-%m') AS month, i.uid, i.id AS income_id, i.total,
    i.created_at AS income_created_at, i.type_income, i.transaction_id, tce.category_id, tce.id AS t_id, tce.title AS category_title
FROM 
    incomes i
LEFT JOIN 
    t_category_incomes tce ON i.id = tce.category_id
    where uid = ? && DATE_FORMAT(i.created_at, '%Y-%m') = ?
ORDER BY month ASC, i.uid ASC; `
	rows, err := tx.QueryContext(ctx, script, param.Uid, param.Month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incomes := []model.MonthlyIDetail{}
	var i model.MonthlyIDetail
	for rows.Next() {
		err := rows.Scan(
			&i.Month, &i.Uid, &i.Id,
			&i.Total, &i.CreatedAt, &i.TypeIncome,
			&i.TransactionId, &i.TCategory.CategoryId,
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

// GetMonthlyExpenseReport implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) GetMonthlyExpenseReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error) {
	script := `select YEAR(created_at) AS year,MONTH(created_at) AS month, uid, sum(total) as total_expenses  
	from expenses where uid = ? GROUP BY YEAR(created_at), MONTH(created_at), uid;`
	rows, err := tx.QueryContext(ctx, script, uid.Uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reports := []model.MonthlyReportResponse{}
	for rows.Next() {
		var i model.MonthlyReport
		err := rows.Scan(
			&i.Year.Float64,
			&i.Month.Float64,
			&i.Uid.String,
			&i.TotalExpense)
		if err != nil {
			return nil, err
		}
		v, er := json.Marshal(&i)
		if er != nil {
			return nil, er
		}
		e := json.Unmarshal(v, &i)
		if e != nil {
			return nil, e
		}
		r := model.MonthlyReportResponse{
			Month:        i.Month.Float64,
			Year:         i.Year.Float64,
			TotalExpense: i.TotalExpense.Float64,
		}
		reports = append(reports, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

// GetMonthlyIncomeReport implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) GetMonthlyIncomeReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error) {
	script := `select YEAR(created_at) AS year, MONTH(created_at) AS month, uid, sum(total) as total_income  
	from incomes where uid = ? GROUP BY YEAR(created_at), MONTH(created_at), uid;`
	rows, err := tx.QueryContext(ctx, script, uid.Uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reports := []model.MonthlyReportResponse{}
	for rows.Next() {
		var i model.MonthlyReport
		err := rows.Scan(
			&i.Year.Float64,
			&i.Month.Float64,
			&i.Uid.String,
			&i.TotalIncome)
		if err != nil {
			return nil, err
		}
		v, er := json.Marshal(&i)
		if er != nil {
			return nil, er
		}
		e := json.Unmarshal(v, &i)
		if e != nil {
			return nil, e
		}
		r := model.MonthlyReportResponse{
			Month:       i.Month.Float64,
			Year:        i.Year.Float64,
			TotalIncome: i.TotalIncome.Float64,
		}
		reports = append(reports, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

// @Depreceted Not used will be remove later MonthlyReportExpensesDetail implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) MonthlyReportExpensesDetail(ctx context.Context, tx *sql.Tx,
	param model.ParamMonthlyReportDetail) ([]model.MonthlyXDetail, error) {
	script := `
SELECT 
    DATE_FORMAT(e.created_at, '%Y-%m') AS month,
    e.uid, e.id AS expense_id, e.expense_type, e.total, e.Notes, e.created_at AS expense_created_at, 
    e.status AS expense_status, e.transaction_id, tce.category_id,tce.id AS t_id, tce.title AS category_title
FROM 
    expenses e
LEFT JOIN 
    t_category_expenses tce ON e.id = tce.category_id
    where uid = ? && status ='success' && DATE_FORMAT(e.created_at, '%Y-%m') = ?
ORDER BY month ASC, e.uid ASC;`
	rows, err := tx.QueryContext(ctx, script, param.Uid, param.Month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := []model.MonthlyXDetail{}
	for rows.Next() {
		var i model.MonthlyXDetail
		err := rows.Scan(
			&i.Month, &i.Uid, &i.Id, &i.ExpenseType, &i.Total, &i.Notes,
			&i.CreatedAt, &i.Status, &i.TransactionId,
			&i.TCategory.CategoryId, &i.TCategory.Id, &i.TCategory.Title,
		)
		if err != nil {
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

func NewMonthlyRepository() MonthlyReportRepository {
	return MonthlyReportRepositoryImpl{}
}
