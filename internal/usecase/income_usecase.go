package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/msarifin29/be_budget_in/util/zero"
	"github.com/sirupsen/logrus"
)

type IncomeUsecase interface {
	CreateIncome(ctx context.Context, params model.CreateIncomeParams) (model.IncomeResponse, error)
	GetIncomes(ctx context.Context, params model.GetIncomeParams) ([]model.IncomeResponse, float64, error)
	GetIncomesByMonth(ctx context.Context, params model.MonthlyParams) (float64, error)
}

type IncomeUsecaseImpl struct {
	CategoryRepo repository.CategoryRepository
	IncomeRepo   repository.IncomeRepository
	// BalanceRepo  repository.BalanceRepository
	AccountRepo repository.AccountRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetIncomesByMonth implements IncomeUsecase.
func (u *IncomeUsecaseImpl) GetIncomesByMonth(ctx context.Context, params model.MonthlyParams) (float64, error) {
	var totalIncomes float64
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	expenses, err := u.IncomeRepo.GetIncomesByMonth(ctx, tx, params)
	if err != nil {
		u.Log.Error(err)
		err = errors.New("failed get incomes")
		return 0, err
	}
	for _, ex := range expenses {
		totalIncomes += ex.Total
	}
	return totalIncomes, nil
}

// CreateIncome implements IncomeUsecase.
func (u *IncomeUsecaseImpl) CreateIncome(ctx context.Context, params model.CreateIncomeParams) (model.IncomeResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	req := model.Income{
		Uid:            params.Uid,
		CategoryIncome: params.CategoryIncome,
		TypeIncome:     params.TypeIncome,
		Total:          params.Total,
		TransactionId:  uuid.NewString(),
		CreatedAt:      util.CreatedAt(params.CreatedAt),
	}
	err := NewIncome(ctx, tx, u.AccountRepo, u.Log, params.TypeIncome, params.AccountId, params.Total)
	if err != nil {
		return model.IncomeResponse{}, err
	}
	res, err := u.IncomeRepo.CreateIncome(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create income %v", err)
		return model.IncomeResponse{}, err
	}
	paramCategory := model.Category{
		CategoryId: res.Id,
		Id:         float64(params.CategoryId),
		Title:      util.InputCategoryIncome(float64(params.CategoryId)),
	}

	category, catErr := u.CategoryRepo.CreateCategoryIncomes(ctx, tx, paramCategory)
	if catErr != nil {
		u.Log.Errorf("failed create catgory income %v", catErr)
		return model.IncomeResponse{}, catErr
	}

	update := zero.TimeFromPtr(res.UpdatedAt)
	return model.IncomeResponse{
		Uid:            res.Uid,
		Id:             res.Id,
		CategoryIncome: category.Title,
		TypeIncome:     res.TypeIncome,
		Total:          req.Total,
		TransactionId:  res.TransactionId,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      &update.Time,
	}, nil
}

// GetIncomes implements IncomeUsecase.
func (u *IncomeUsecaseImpl) GetIncomes(ctx context.Context, params model.GetIncomeParams) ([]model.IncomeResponse, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	total, err := u.IncomeRepo.GetTotalIncomes(ctx, tx, params.Uid, params.CategoryId, params.TypeIncome)
	if err != nil {
		u.Log.Errorf("failed get total incomes %v ", err)
		return []model.IncomeResponse{}, 0, err
	}
	incomes, err := u.IncomeRepo.GetIncomes(ctx, tx, params)
	if err != nil {
		u.Log.Errorf("failed get incomes %v ", err)
		return []model.IncomeResponse{}, 0, err
	}
	return incomes, total, nil
}

func NewIncomeUsecase(IncomeRepo repository.IncomeRepository,
	// BalanceRepo repository.BalanceRepository,
	AccountRepo repository.AccountRepository,
	Log *logrus.Logger, db *sql.DB,
	CategoryRepo repository.CategoryRepository) IncomeUsecase {
	return &IncomeUsecaseImpl{IncomeRepo: IncomeRepo,
		// BalanceRepo: BalanceRepo,
		AccountRepo: AccountRepo, Log: Log,
		db: db, CategoryRepo: CategoryRepo}
}

func NewIncome(ctx context.Context, tx *sql.Tx, accountRepo repository.AccountRepository, Log *logrus.Logger, typeIncome string, accountId string, input float64) error {
	var newCash, newDebit float64

	switch typeIncome {
	case util.CASH:
		account, err := accountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
		cash := account.Cash
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
		err = accountRepo.UpdateAccountCash(ctx, tx, model.UpdateAccountCash{AccountId: accountId, Cash: newCash})
		if err != nil {
			err = errors.New("failed update cash")
			Log.Error(err)
			return err
		}
	case util.DEBIT:
		account, err := accountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
		debit := account.Balance
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
		err = accountRepo.UpdateAccountBalance(ctx, tx, model.UpdateAccountBalance{AccountId: accountId, Balance: newDebit})
		if err != nil {
			err = errors.New("failed update debit")
			Log.Error(err)
			return err
		}
	}
	return nil
}
