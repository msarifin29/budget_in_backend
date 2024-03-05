package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/msarifin29/be_budget_in/util/zero"
	"github.com/sirupsen/logrus"
)

type CreditUsecase interface {
	CreateCredit(ctx context.Context, params model.CreateCreditRequest) (model.Credit, error)
	GetAllCredit(ctx context.Context, params model.GetCreditParams) ([]model.Credit, float64, error)
	GetAllHistoryCredit(ctx context.Context, params model.GetHistoriesCreditParams) ([]model.HistoryCredit, float64, error)
	UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditParams) (model.UpdateHistoryResponse, error)
}

type CreditUsecaseImpl struct {
	CreditRepo  repository.CreditRepository
	BalanceRepo repository.BalanceRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetAllHistoryCredit implements CreditUsecase.
func (u *CreditUsecaseImpl) GetAllHistoryCredit(ctx context.Context, params model.GetHistoriesCreditParams) ([]model.HistoryCredit, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	historiesCredits, err := u.CreditRepo.GetAllHistoryCredit(ctx, tx, params)
	if err != nil {
		u.Log.Error()
		err = errors.New("failed get all history credit")
		return []model.HistoryCredit{}, 0, err
	}
	count, err := u.CreditRepo.GetCountHistoryCredit(ctx, tx, params.CreditId)
	if err != nil {
		u.Log.Error()
		err = errors.New("failed get count history credit")
		return []model.HistoryCredit{}, 0, err
	}
	return historiesCredits, count, nil
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
func (u *CreditUsecaseImpl) GetAllCredit(ctx context.Context, params model.GetCreditParams) ([]model.Credit, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	credits, err := u.CreditRepo.GetAllCredit(ctx, tx, params)
	if err != nil {
		u.Log.Error()
		err = errors.New("failed get all credit")
		return []model.Credit{}, 0, err
	}
	count, err := u.CreditRepo.GetCountCredit(ctx, tx, params.Uid)
	if err != nil {
		u.Log.Error()
		err = errors.New("failed get count credit")
		return []model.Credit{}, 0, err
	}
	return credits, count, nil
}

// UpdateHistoryCredit implements CreditUsecase.
func (u *CreditUsecaseImpl) UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditParams) (model.UpdateHistoryResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	credit, err := u.CreditRepo.GetCreditById(ctx, tx, model.GetCreditRequest{Uid: params.Uid, Id: params.CreditId})
	if err != nil {
		u.Log.Errorf("failed get credit with credit id %v", params.CreditId)
		return model.UpdateHistoryResponse{}, err
	}
	if params.Status == util.ACTIVE {
		err := errors.New("cannot update history credit with same value")
		return model.UpdateHistoryResponse{}, err
	}

	reqId := model.GetHistoryCreditRequest{Uid: params.Uid, Id: params.Id}
	historyC, err := u.CreditRepo.GetHistoryCreditById(ctx, tx, reqId)
	if err != nil {
		u.Log.Errorf("failed get history credit with id %v", params.Id)
		return model.UpdateHistoryResponse{}, err
	}
	if historyC.Status != util.ACTIVE {
		err = errors.New("status credit is completed")
		u.Log.Error(err)
		return model.UpdateHistoryResponse{}, err
	}

	ok, err := u.CreditRepo.UpdateHistoryCredit(ctx, tx, params)
	if err != nil || !ok {
		u.Log.Error(err)
		err = errors.New("failed update history credit")
		return model.UpdateHistoryResponse{}, err
	}
	err = UpdateTotalBalanceOrCash(ctx, tx, u.CreditRepo, u.BalanceRepo, historyC, params.TypePayment, params.Uid, u.Log)
	if err != nil {
		return model.UpdateHistoryResponse{}, err
	}
	newCredit := credit.Total - historyC.Total
	ok, totalErr := u.CreditRepo.UpdateTotalCredit(ctx, tx, params.Uid, params.CreditId, newCredit)
	if totalErr != nil || !ok {
		u.Log.Error(err)
		err = errors.New("failed update total credit")
		return model.UpdateHistoryResponse{}, err
	}

	if newCredit <= 0 {
		ok, err := u.CreditRepo.UpdateCredit(ctx, tx, model.UpdateCreditRequest{Uid: params.Uid, Id: params.CreditId, StatusCredit: util.COMPLETED})
		if !ok || err != nil {
			u.Log.Error(err)
			err = errors.New("failed update credit")
			return model.UpdateHistoryResponse{}, err
		}
		// Update debts from user
		newCreditCompletted := credit.Installment * credit.LoanTerm
		err = NewDebts(ctx, tx, u.Log, u.BalanceRepo, util.COMPLETED, params.Uid, newCreditCompletted)
		if err != nil {
			u.Log.Error(err)
			err = errors.New("failed update debts user")
			return model.UpdateHistoryResponse{}, err
		}
	}

	return model.UpdateHistoryResponse{
		Id:          historyC.Id,
		Th:          historyC.Th,
		Total:       historyC.Total,
		Status:      params.Status,
		TypePayment: params.TypePayment,
		CreatedAt:   time.Now(),
	}, nil
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

func UpdateTotalBalanceOrCash(ctx context.Context, tx *sql.Tx,
	creditRepo repository.CreditRepository,
	BalanceRepo repository.BalanceRepository,
	historyCredit model.HistoryCredit,
	typePayment string,
	uid string,
	Log *logrus.Logger,
) error {
	switch typePayment {
	case util.DEBIT:
		balance, err := BalanceRepo.GetBalance(ctx, tx, uid)
		if err != nil {
			err = fmt.Errorf("failed get balance from uid %v", uid)
			return err
		}
		if balance < historyCredit.Total {
			err = errors.New("cannot upgrade balance with total greater than balance")
			return err
		}
		newBalance := balance - historyCredit.Total
		Log.Infof("newbalance = %v, balance = %v, input = %v", newBalance, balance, historyCredit.Total)
		err = BalanceRepo.SetBalance(ctx, tx, uid, newBalance)
		if err != nil {
			err = errors.New("failed update balance")
			return err
		}
	case util.CASH:
		cash, err := BalanceRepo.GetCash(ctx, tx, uid)
		if err != nil {
			err = fmt.Errorf("failed get cash from uid %v", uid)
			return err
		}
		if cash < historyCredit.Total {
			err = errors.New("cannot upgrade cash with total greater than cash")
			return err
		}
		newCash := cash - historyCredit.Total
		Log.Infof("newCash = %v, cash = %v, input = %v", newCash, cash, historyCredit.Total)
		err = BalanceRepo.SetCash(ctx, tx, uid, newCash)
		if err != nil {
			err = errors.New("failed update cash")
			return err
		}
	}
	return nil
}
