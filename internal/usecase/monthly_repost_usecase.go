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
	return res, nil
}

func NewMonthlyReportUsecase(MonthlyRepo repository.MonthlyReportRepository, Log *logrus.Logger, db *sql.DB) MonthlyReportUsecase {
	return &MonthlyReportUsecaseImpl{MonthlyRepo: MonthlyRepo, Log: Log, db: db}
}
