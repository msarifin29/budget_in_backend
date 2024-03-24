package util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/sirupsen/logrus"
)

func NewBalance(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, accountRepo repository.AccountRepository, status string, accountId string, input float64) error {
	var newBalance float64
	account, err := accountRepo.GetAccountByUserId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
	balance := account.Balance
	if err != nil {
		err = errors.New("failed get balance")
		Log.Error(err)
		return err
	}

	switch status {
	case SUCCESS:
		if balance < input {
			err = fmt.Errorf("invalid input, your balance is %v ", balance)
			Log.Error(err)
			return err
		}
		newBalance = balance - input
	case CANCELLED:
		newBalance = balance + input
	}

	Log.Infof("newbalance = %v, balance = %v, input = %v", newBalance, balance, input)
	err = accountRepo.UpdateAccountBalance(ctx, tx, model.UpdateAccountBalance{AccountId: accountId, Balance: newBalance})
	if err != nil {
		err = errors.New("failed update balance")
		Log.Error(err)
		return err
	}
	return nil
}

func NewCash(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, accountRepo repository.AccountRepository, status string, accountId string, input float64) error {
	var newCash float64
	account, err := accountRepo.GetAccountByUserId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
	cash := account.Cash
	if err != nil {
		err = errors.New("failed get cash")
		Log.Error(err)
		return err
	}

	switch status {
	case SUCCESS:
		if cash < input {
			err = fmt.Errorf("invalid input, your cash is %v ", cash)
			Log.Error(err)
			return err
		}
		newCash = cash - input
	case CANCELLED:
		newCash = cash + input
	}
	Log.Infof("newCash = %v, cash = %v, input = %v", newCash, cash, input)
	err = accountRepo.UpdateAccountCash(ctx, tx, model.UpdateAccountCash{AccountId: accountId, Cash: newCash})
	if err != nil {
		err = errors.New("failed update cash")
		Log.Error(err)
		return err
	}
	return nil
}

func NewDebts(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, accountRepo repository.AccountRepository, status string, accountId string, input float64) error {
	var newDebt float64
	account, err := accountRepo.GetAccountByUserId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
	debts := account.Debts
	if err != nil {
		err = errors.New("failed get debts")
		Log.Error(err)
		return err
	}
	switch status {
	case ACTIVE:
		if input <= 0 {
			err = fmt.Errorf("invalid input, cannot add debts with value %v ", input)
			Log.Error(err)
			return err
		}
		newDebt = debts + input
	case COMPLETED:
		newDebt = debts - input
	}

	Log.Infof("newDebt = %v, debts = %v, input = %v", newDebt, debts, input)
	err = accountRepo.UpdateAccountDebts(ctx, tx, model.UpdateAccountDebts{AccountId: accountId, Debts: newDebt})
	if err != nil {
		err = errors.New("failed update debts")
		Log.Error(err)
		return err
	}
	return nil
}
