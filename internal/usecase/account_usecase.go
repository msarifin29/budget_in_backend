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

type AccountUsacase interface {
	CreateAccount(ctx context.Context, account model.CreateAccountRequest) (model.Account, error)
	GetAllAccounts(ctx context.Context, param model.GetAllAccountRequest) ([]model.Account, error)
	UpdateMaxBudget(ctx context.Context, account model.UpdateMaxBudgetRequest) (bool, error)
	GetMaxBudget(ctx context.Context, account model.GetMaxBudgetParam) (model.MaxBudgetResponse, error)
}

type AccountUsacaseImpl struct {
	AccountRepo repository.AccountRepository
	ExpenseRepo repository.ExpenseRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetAllAccounts implements AccountUsacase.
func (u *AccountUsacaseImpl) GetAllAccounts(ctx context.Context, param model.GetAllAccountRequest) ([]model.Account, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	accounts, err := u.AccountRepo.GetAllAccount(ctx, tx, param)
	if err != nil {
		u.Log.Errorf("Failed get max budget %e", err)
		err = errors.New("failed get max budget")
		return []model.Account{}, err
	}

	return accounts, nil
}

// GetMaxBudget implements AccountUsacase.
func (u *AccountUsacaseImpl) GetMaxBudget(ctx context.Context, account model.GetMaxBudgetParam) (model.MaxBudgetResponse, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	totalExpense := 0
	resAccount, er := u.AccountRepo.GetAccountByUserId(ctx, tx, account.Uid)
	if er != nil {
		u.Log.Errorf("Failed get max budget %e", er)
		er = errors.New("failed get max budget")
		return model.MaxBudgetResponse{}, er
	}
	totalEx, err := u.ExpenseRepo.GetExpenseThisMonth(ctx, tx, account.Uid)
	if err != nil {
		u.Log.Errorf("Failed get expense %e", err)
		err = errors.New("failed get expense")
		return model.MaxBudgetResponse{}, err
	}
	totalExpense = int(totalEx)

	return model.MaxBudgetResponse{
		Uid:          resAccount.UserId,
		AccountId:    resAccount.AccountId,
		TotalExpense: float64(totalExpense),
		MaxBudget:    resAccount.MaxBudget}, nil
}

// UpdateMaxBudget implements AccountUsacase.
func (u *AccountUsacaseImpl) UpdateMaxBudget(ctx context.Context, account model.UpdateMaxBudgetRequest) (bool, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	err := u.AccountRepo.UpdateMaxBudget(ctx, tx, account)
	if err != nil {
		u.Log.Errorf("Failed update max budget %e", err)
		err = errors.New("failed update max budget")
		return false, err
	}
	return true, nil
}

// CreateAccount implements AccountUsacase.
func (u *AccountUsacaseImpl) CreateAccount(ctx context.Context, account model.CreateAccountRequest) (model.Account, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	req := model.Account{
		UserId:      account.UserId,
		AccountId:   uuid.NewString(),
		AccountName: account.AccountName,
		Balance:     account.Balance,
		Cash:        account.Cash,
		Debts:       0,
		Savings:     0,
		Currency:    account.Currency,
	}
	res, err := u.AccountRepo.CreateAccount(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create new account :%v", err)
		err = errors.New("failed create new account")
		return model.Account{}, err
	}
	update := zero.TimeFromPtr(res.UpdatedAt)
	return model.Account{
		UserId:      res.UserId,
		AccountId:   res.AccountId,
		AccountName: res.AccountName,
		Balance:     res.Balance,
		Cash:        res.Cash,
		Debts:       res.Debts,
		Savings:     res.Savings,
		Currency:    res.Currency,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   &update.Time,
	}, nil
}

func NewAccountUsacase(AccountRepo repository.AccountRepository, ExpenseRepo repository.ExpenseRepository, Log *logrus.Logger, db *sql.DB) AccountUsacase {
	return &AccountUsacaseImpl{AccountRepo: AccountRepo, ExpenseRepo: ExpenseRepo, Log: Log, db: db}
}
