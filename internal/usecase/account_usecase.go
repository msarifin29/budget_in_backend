package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/msarifin29/be_budget_in/util/zero"
	"github.com/sirupsen/logrus"
)

type AccountUsacase interface {
	CreateAccount(ctx context.Context, account model.CreateAccountRequest) (model.Account, error)
}

type AccountUsacaseImpl struct {
	AccountRepo repository.AccountRepository
	Log         *logrus.Logger
	db          *sql.DB
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
	update := zero.TimeFromPtr(&res.UpdatedAt)
	return model.Account{
		UserId:      res.UserId,
		AccountId:   res.AccountId,
		AccountName: res.AccountName,
		Balance:     res.Balance,
		Cash:        res.Cash,
		Debts:       res.Debts,
		Savings:     res.Savings,
		Currency:    res.Currency,
		CreatedAt:   time.Now(),
		UpdatedAt:   update.Time,
	}, nil
}

func NewAccountUsacase(AccountRepo repository.AccountRepository, Log *logrus.Logger, db *sql.DB) AccountUsacase {
	return &AccountUsacaseImpl{AccountRepo: AccountRepo, Log: Log, db: db}
}
