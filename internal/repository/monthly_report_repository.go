package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type MonthlyReportRepository interface {
	GetMonthlyReport(ctx context.Context, tx *sql.Tx, uid string) ([]model.MonthlyReportResponse, error)
}

type MonthlyReportRepositoryImpl struct{}

// GetMonthlyReport implements MonthlyReportRepository.
func (MonthlyReportRepositoryImpl) GetMonthlyReport(ctx context.Context, tx *sql.Tx, uid string) ([]model.MonthlyReportResponse, error) {
	script := `SELECT
YEAR(e.created_at) AS year,
MONTH(e.created_at) AS month,
e.uid,
SUM(e.total) AS total_expense,
SUM(i.total) AS total_income
FROM expenses AS e
LEFT JOIN incomes AS i ON e.uid = i.uid AND YEAR(e.created_at) = YEAR(i.created_at) AND MONTH(e.created_at) = MONTH(i.created_at)
WHERE e.uid = ? 
GROUP BY YEAR(e.created_at), MONTH(e.created_at), e.uid
ORDER BY YEAR(e.created_at), MONTH(e.created_at);`
	rows, err := tx.QueryContext(ctx, script, uid)
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
