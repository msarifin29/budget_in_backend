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
	GetMonthlyReport(ctx context.Context, params model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error)
	GetMonthlyReportDetail(ctx context.Context, params model.ParamMonthlyReportDetail) (model.MonthlyReportDetailResponse, error)
}
type MonthlyReportUsecaseImpl struct {
	MonthlyRepo repository.MonthlyReportRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetMonthlyReportDetail implements MonthlyReportUsecase.
func (u *MonthlyReportUsecaseImpl) GetMonthlyReportDetail(ctx context.Context, params model.ParamMonthlyReportDetail) (model.MonthlyReportDetailResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	expenses, errX := u.MonthlyRepo.MonthlyReportExpensesDetail(ctx, tx, params)
	if errX != nil {
		u.Log.Errorf("failed get monthly expenses detail %v", errX)
		errX = errors.New("failed get monthly expenses detail")
		return model.MonthlyReportDetailResponse{}, errX
	}
	incomes, errI := u.MonthlyRepo.MonthlyReportIncomesDetail(ctx, tx, params)
	if errI != nil {
		u.Log.Errorf("failed get monthly incomes detail %v", errI)
		errI = errors.New("failed get monthly incomes detail")
		return model.MonthlyReportDetailResponse{}, errI
	}
	return model.MonthlyReportDetailResponse{
		ExpensesRecords: expenses,
		IncomesRecords:  incomes,
	}, nil
}

// GetMonthlyReport implements MonthlyReportUsecase.
func (u *MonthlyReportUsecaseImpl) GetMonthlyReport(ctx context.Context, params model.ParamMonthlyReport) ([]model.MonthlyReportResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	in, inErr := u.MonthlyRepo.GetMonthlyIncomeReport(ctx, tx, params)
	if inErr != nil {
		u.Log.Errorf("failed get monthly income report %v", inErr)
		inErr = errors.New("failed get monthly income report")
		return []model.MonthlyReportResponse{}, inErr
	}
	ex, exErr := u.MonthlyRepo.GetMonthlyExpenseReport(ctx, tx, params)
	if exErr != nil {
		u.Log.Errorf("failed get monthly expenses report %v", exErr)
		exErr = errors.New("failed get monthly expenses report")
		return []model.MonthlyReportResponse{}, exErr
	}

	reports := []model.MonthlyReportResponse{}
	for i := 1; i < 13; i++ {
		var m model.MonthlyReportResponse
		for _, v := range in {
			m.Year = v.Year
			m.Month = float64(i)

			if len(in) != 0 {
				if int(v.Month) == i {
					m.TotalIncome = v.TotalIncome
				}
			} else {
				m.TotalIncome = 0
			}
		}
		for _, expen := range ex {
			m.Year = expen.Year
			m.Month = float64(i)

			if (len(ex)) != 0 {
				if int(expen.Month) == i {
					m.TotalExpense = expen.TotalExpense
				}
			} else {
				m.TotalExpense = 0
			}
		}
		reports = append(reports, m)
	}
	return reports, nil
}

func NewMonthlyReportUsecase(MonthlyRepo repository.MonthlyReportRepository, Log *logrus.Logger, db *sql.DB) MonthlyReportUsecase {
	return &MonthlyReportUsecaseImpl{MonthlyRepo: MonthlyRepo, Log: Log, db: db}
}
