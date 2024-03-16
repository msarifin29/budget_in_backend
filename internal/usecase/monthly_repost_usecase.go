package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type MonthlyReportUsecase interface {
	GetMonthlyReport(ctx context.Context, uid string) ([]model.MonthlyReportResponse, error)
}
type MonthlyReportUsecaseImpl struct {
	MonthlyRepo repository.MonthlyReportRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetMonthlyReport implements MonthlyReportUsecase.
func (u *MonthlyReportUsecaseImpl) GetMonthlyReport(ctx context.Context, uid string) ([]model.MonthlyReportResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	res, err := u.MonthlyRepo.GetMonthlyReport(ctx, tx, uid)
	if err != nil {
		u.Log.Errorf("failed get monthly report %v", err)
		err = errors.New("failed get monthly report")
		return []model.MonthlyReportResponse{}, err
	}
	reports := []model.MonthlyReportResponse{}
	for i := 1; i < 13; i++ {
		var m model.MonthlyReportResponse
		for _, v := range res {
			if len(res) != 0 {
				m.Year = v.Year
				m.Month = float64(i)
				if int(v.Month) == i {
					m.TotalExpense = v.TotalExpense
					m.TotalIncome = v.TotalIncome
				}
			} else {
				m.Year = 0
				m.Month = 0
				m.TotalExpense = 0
				m.TotalIncome = 0
			}

		}
		reports = append(reports, m)

	}
	return reports, nil
}

func NewMonthlyReportUsecase(MonthlyRepo repository.MonthlyReportRepository, Log *logrus.Logger, db *sql.DB) MonthlyReportUsecase {
	return &MonthlyReportUsecaseImpl{MonthlyRepo: MonthlyRepo, Log: Log, db: db}
}
