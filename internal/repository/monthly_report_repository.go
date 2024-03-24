package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type MonthlyReportRepository interface {
	GetMonthlyReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
	GetMonthlyIncomeReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
	GetMonthlyExpenseReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
}

type MonthlyReportRepositoryImpl struct{}

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

// @Depreceted Not used will be remove later GetMonthlyReport implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) GetMonthlyReport(ctx context.Context, tx *sql.Tx, uid model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error) {
	script := `WITH income_and_expense AS (
		SELECT
		  YEAR(i.created_at) AS year,
		  MONTH(i.created_at) AS month,
		  i.uid,
		  SUM(i.total) AS total_income,
		  SUM(CASE WHEN e.status = 'success' AND YEAR(e.created_at) = YEAR(i.created_at) 
		  AND MONTH(e.created_at) = MONTH(i.created_at) 
		  AND DAY(e.created_at) = DAY(i.created_at) THEN e.total ELSE 0 END) AS total_expense
		FROM incomes AS i
		LEFT JOIN expenses AS e ON i.uid = e.uid
		WHERE i.uid = ?  
		GROUP BY YEAR(i.created_at), MONTH(i.created_at), i.uid
	  )
	  SELECT * FROM income_and_expense;
	  `
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
			&i.TotalExpense.Float64,
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
			Month:        i.Month.Float64,
			Year:         i.Year.Float64,
			TotalExpense: i.TotalExpense.Float64,
			TotalIncome:  i.TotalIncome.Float64,
		}
		reports = append(reports, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func NewMonthlyRepository() MonthlyReportRepository {
	return MonthlyReportRepositoryImpl{}
}
