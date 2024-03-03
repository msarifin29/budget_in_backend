package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/msarifin29/be_budget_in/util/zero"
	"github.com/sirupsen/logrus"
)

type IncomeUsecase interface {
	CreateIncome(ctx context.Context, params model.CreateIncomeRequest) (model.IncomeResponse, error)
	GetIncomes(ctx context.Context, params model.GetIncomeParams) ([]model.Income, float64, error)
}

type IncomeUsecaseImpl struct {
	IncomeRepo  repository.IncomeRepository
	BalanceRepo repository.BalanceRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// CreateIncome implements IncomeUsecase.
func (u *IncomeUsecaseImpl) CreateIncome(ctx context.Context, params model.CreateIncomeRequest) (model.IncomeResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	req := model.Income{
		Uid:            params.Uid,
		CategoryIncome: params.CategoryIncome,
		TypeIncome:     params.TypeIncome,
		Total:          params.Total,
	}
	err := NewIncome(ctx, tx, u.BalanceRepo, u.Log, params.TypeIncome, params.Uid, params.Total)
	if err != nil {
		return model.IncomeResponse{}, err
	}
	res, err := u.IncomeRepo.CreateIncome(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create income %v", err)
		return model.IncomeResponse{}, err
	}
	update := zero.TimeFromPtr(&res.UpdatedAt)
	return model.IncomeResponse{
		Uid:            res.Uid,
		Id:             res.Id,
		CategoryIncome: res.CategoryIncome,
		TypeIncome:     res.TypeIncome,
		Total:          req.Total,
		CreatedAt:      time.Now(),
		UpdatedAt:      update.Time,
	}, nil
}

// GetIncomes implements IncomeUsecase.
func (u *IncomeUsecaseImpl) GetIncomes(ctx context.Context, params model.GetIncomeParams) ([]model.Income, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	total, err := u.IncomeRepo.GetTotalIncomes(ctx, tx, params.Uid, params.CategoryIncome)
	if err != nil {
		u.Log.Errorf("failed get total incomes %v ", err)
		return []model.Income{}, 0, err
	}
	incomes, err := u.IncomeRepo.GetIncomes(ctx, tx, params)
	if err != nil {
		u.Log.Errorf("failed get incomes %v ", err)
		return []model.Income{}, 0, err
	}
	return incomes, total, nil
}

func NewIncomeUsecase(IncomeRepo repository.IncomeRepository, BalanceRepo repository.BalanceRepository, Log *logrus.Logger, db *sql.DB) IncomeUsecase {
	return &IncomeUsecaseImpl{IncomeRepo: IncomeRepo, BalanceRepo: BalanceRepo, Log: Log, db: db}
}

func NewIncome(ctx context.Context, tx *sql.Tx, balanceRepo repository.BalanceRepository, Log *logrus.Logger, typeIncome string, uid string, input float64) error {
	var newCash, newDebit float64

	switch typeIncome {
	case util.CASH:
		cash, err := balanceRepo.GetCash(ctx, tx, uid)
		if err != nil {
			err = errors.New("failed get cash")
			Log.Error(err)
			return err
		}
		if input <= 0 {
			err = errors.New("invalid input, min income 2000")
			Log.Error(err)
			return err
		}
		newCash = cash + input
		Log.Infof("newCash = %v, cash = %v, input = %v", newCash, cash, input)
		err = balanceRepo.SetCash(ctx, tx, uid, newCash)
		if err != nil {
			err = errors.New("failed update cash")
			Log.Error(err)
			return err
		}
	case util.DEBIT:
		debit, err := balanceRepo.GetBalance(ctx, tx, uid)
		if err != nil {
			err = errors.New("failed get debit")
			Log.Error(err)
			return err
		}
		if input <= 0 {
			err = errors.New("invalid input, min income 2000")
			Log.Error(err)
			return err
		}
		newDebit = debit + input
		Log.Infof("newdebit = %v, debit = %v, input = %v", newDebit, debit, input)
		err = balanceRepo.SetBalance(ctx, tx, uid, newDebit)
		if err != nil {
			err = errors.New("failed update debit")
			Log.Error(err)
			return err
		}
	}
	return nil
}
