package usecase

import (
	"context"
	"database/sql"
	"errors"

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
	CashWithdrawal(ctx context.Context, params model.CashWithdrawalParam) (bool, error)
	TopUp(ctx context.Context, params model.TopUpParam) (bool, error)
}

type IncomeUsecaseImpl struct {
	CategoryRepo repository.CategoryRepository
	IncomeRepo   repository.IncomeRepository
	// BalanceRepo  repository.BalanceRepository
	AccountRepo repository.AccountRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// TopUp implements IncomeUsecase.
func (u *IncomeUsecaseImpl) TopUp(ctx context.Context, params model.TopUpParam) (bool, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	account, errA := u.AccountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: params.AccountId})
	if errA != nil {
		u.Log.Errorf("failed get account %s", errA)
		errA = errors.New("failed get account")
		return false, errA
	} else if account.Cash < params.Total {
		u.Log.Errorf("the cash is not sufficient %s", errA)
		errA = errors.New("the cash is not sufficient")
		return false, errA
	}
	newCash := account.Cash - params.Total
	errC := u.AccountRepo.UpdateAccountCash(ctx, tx, model.UpdateAccountCash{AccountId: params.AccountId, Cash: newCash})
	if errC != nil {
		u.Log.Errorf("failed update cash %s", errC)
		errC = errors.New("failed update cash")
		return false, errC
	}
	newBalance := account.Balance + params.Total
	errU := u.AccountRepo.UpdateAccountBalance(ctx, tx, model.UpdateAccountBalance{AccountId: params.AccountId, Balance: newBalance})
	if errU != nil {
		u.Log.Errorf("failed update balance %s", errU)
		errU = errors.New("failed update balance")
		return false, errU
	}
	return true, nil
}

// CashWithdrawal implements IncomeUsecase.
func (u *IncomeUsecaseImpl) CashWithdrawal(ctx context.Context, params model.CashWithdrawalParam) (bool, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	account, errA := u.AccountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: params.AccountId})
	if errA != nil {
		u.Log.Errorf("failed get account %s", errA)
		errA = errors.New("failed get account")
		return false, errA
	} else if account.Balance < params.Total {
		u.Log.Errorf("the balance is not sufficient %s", errA)
		errA = errors.New("the balance is not sufficient")
		return false, errA
	}
	newBlance := account.Balance - params.Total
	errU := u.AccountRepo.UpdateAccountBalance(ctx, tx, model.UpdateAccountBalance{AccountId: params.AccountId, Balance: newBlance})
	if errU != nil {
		u.Log.Errorf("failed update balance %s", errU)
		errU = errors.New("failed update balance")
		return false, errU
	}

	newCash := account.Cash + params.Total
	errC := u.AccountRepo.UpdateAccountCash(ctx, tx, model.UpdateAccountCash{AccountId: params.AccountId, Cash: newCash})
	if errC != nil {
		u.Log.Errorf("failed update cash %s", errC)
		errC = errors.New("failed update cash")
		return false, errC
	}
	return true, nil
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
		TransactionId:  util.RandomString(10),
		CreatedAt:      util.CreatedAt(params.CreatedAt),
		AccountId:      params.AccountId,
		BankName:       params.BankName,
		BankId:         params.BankId,
		Cid:            params.CategoryId,
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
		UserId:     params.Uid,
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
	AccountRepo repository.AccountRepository,
	Log *logrus.Logger, db *sql.DB,
	CategoryRepo repository.CategoryRepository) IncomeUsecase {
	return &IncomeUsecaseImpl{IncomeRepo: IncomeRepo,
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
