package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/msarifin29/be_budget_in/util/zero"
	"github.com/sirupsen/logrus"
)

type CreditUsecase interface {
	CreateCredit(ctx context.Context, params model.CreateCreditRequest) (model.Credit, error)
	UpdateCredit(ctx context.Context, params model.UpdateCreditRequest) (bool, error)
	UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditRequest) (bool, error)
}

type CreditUsecaseImpl struct {
	CreditRepo  repository.CreditRepository
	BalanceRepo repository.BalanceRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// CreateCredit implements CreditUsecase.
func (u *CreditUsecaseImpl) CreateCredit(ctx context.Context, params model.CreateCreditRequest) (model.Credit, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	req := model.Credit{
		Uid:            params.Uid,
		CategoryCredit: params.CategoryCredit,
		TypeCredit:     params.TypeCredit,
		Total:          params.LoanTerm * params.Installment,
		LoanTerm:       params.LoanTerm,
		StatusCredit:   util.ACTIVE,
		Installment:    params.Installment,
		PaymentTime:    params.PaymentTime,
	}
	creditRes, err := u.CreditRepo.CreateCredit(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create credit %v", err)
		return model.Credit{}, err
	}
	err = NewHistoryCredit(ctx, tx, u.CreditRepo, creditRes)
	if err != nil {
		u.Log.Errorf("failed create history credit %v", err)
		return model.Credit{}, err
	}
	err = NewDebts(ctx, tx, u.Log, u.BalanceRepo, util.ACTIVE, creditRes.Uid, creditRes.Total)
	if err != nil {
		u.Log.Errorf("failed update depts %v", err)
		return model.Credit{}, err
	}
	update := zero.TimeFromPtr(&creditRes.UpdatedAt)
	res := model.Credit{
		Uid:            creditRes.Uid,
		Id:             creditRes.Id,
		CategoryCredit: creditRes.CategoryCredit,
		TypeCredit:     creditRes.TypeCredit,
		Total:          creditRes.Total,
		LoanTerm:       creditRes.LoanTerm,
		StatusCredit:   util.ACTIVE,
		Installment:    creditRes.Installment,
		CreatedAt:      time.Now(),
		UpdatedAt:      update.Time,
		PaymentTime:    creditRes.PaymentTime,
	}
	return res, nil
}

// UpdateCredit implements CreditUsecase.
func (u *CreditUsecaseImpl) UpdateCredit(ctx context.Context, params model.UpdateCreditRequest) (bool, error) {
	panic("unimplemented")
}

// UpdateHistoryCredit implements CreditUsecase.
func (u *CreditUsecaseImpl) UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditRequest) (bool, error) {
	panic("unimplemented")
}

func NewCreditUsecase(CreditRepo repository.CreditRepository, BalanceRepo repository.BalanceRepository, Log *logrus.Logger, db *sql.DB) CreditUsecase {
	return &CreditUsecaseImpl{CreditRepo: CreditRepo, BalanceRepo: BalanceRepo, Log: Log, db: db}
}

func NewHistoryCredit(ctx context.Context, tx *sql.Tx, creditRepo repository.CreditRepository, credit model.Credit) error {
	for i := 0; i < int(credit.LoanTerm); i++ {
		req := model.HistoryCredit{
			CreditId:    credit.Id,
			Th:          float64(i + 1),
			Total:       credit.Installment,
			Status:      util.ACTIVE,
			TypePayment: "",
			PaymentTime: credit.PaymentTime,
		}
		_, err := creditRepo.CreateHistoryCredit(ctx, tx, req)
		if err != nil {
			return err
		}
	}
	return nil
}
