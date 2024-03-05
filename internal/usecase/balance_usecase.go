package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

func NewBalance(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, balanceRepo repository.BalanceRepository, status string, uid string, input float64) error {
	var newBalance float64
	balance, err := balanceRepo.GetBalance(ctx, tx, uid)
	if err != nil {
		err = errors.New("failed get balance")
		Log.Error(err)
		return err
	}

	switch status {
	case util.SUCCESS:
		if balance < input {
			err = fmt.Errorf("invalid input, your balance is %v ", balance)
			Log.Error(err)
			return err
		}
		newBalance = balance - input
	case util.CANCELLED:
		newBalance = balance + input
	}

	Log.Infof("newbalance = %v, balance = %v, input = %v", newBalance, balance, input)
	err = balanceRepo.SetBalance(ctx, tx, uid, newBalance)
	if err != nil {
		err = errors.New("failed update balance")
		Log.Error(err)
		return err
	}
	return nil
}

func NewCash(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, balanceRepo repository.BalanceRepository, status string, uid string, input float64) error {
	var newCash float64
	cash, err := balanceRepo.GetCash(ctx, tx, uid)
	if err != nil {
		err = errors.New("failed get cash")
		Log.Error(err)
		return err
	}

	switch status {
	case util.SUCCESS:
		if cash < input {
			err = fmt.Errorf("invalid input, your cash is %v ", cash)
			Log.Error(err)
			return err
		}
		newCash = cash - input
	case util.CANCELLED:
		newCash = cash + input
	}
	Log.Infof("newCash = %v, cash = %v, input = %v", newCash, cash, input)
	err = balanceRepo.SetCash(ctx, tx, uid, newCash)
	if err != nil {
		err = errors.New("failed update cash")
		Log.Error(err)
		return err
	}
	return nil
}

func NewSavings(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, balanceRepo repository.BalanceRepository, uid string, input float64) error {
	saving, err := balanceRepo.GetSaving(ctx, tx, uid)
	if err != nil {
		err = errors.New("failed get savings")
		Log.Error(err)
		return err
	}
	if input <= 0 {
		err = fmt.Errorf("invalid input, cannot add savings with value %v ", saving)
		Log.Error(err)
		return err
	}
	newSaving := saving + input
	Log.Infof("newSaving = %v, savings = %v, input = %v", newSaving, saving, input)
	err = balanceRepo.SetSaving(ctx, tx, uid, newSaving)
	if err != nil {
		err = errors.New("failed update savings")
		Log.Error(err)
		return err
	}
	return nil
}

func NewDebts(ctx context.Context, tx *sql.Tx, Log *logrus.Logger, balanceRepo repository.BalanceRepository, status string, uid string, input float64) error {
	var newDebt float64
	debts, err := balanceRepo.GetDebt(ctx, tx, uid)
	if err != nil {
		err = errors.New("failed get debts")
		Log.Error(err)
		return err
	}
	switch status {
	case util.ACTIVE:
		if input <= 0 {
			err = fmt.Errorf("invalid input, cannot add debts with value %v ", input)
			Log.Error(err)
			return err
		}
		newDebt = debts + input
	case util.COMPLETED:
		if debts <= 0 {
			err = fmt.Errorf("min debts is 0 %v ", debts)
			Log.Error(err)
			return err
		}
		newDebt = debts - input
	}

	Log.Infof("newDebt = %v, debts = %v, input = %v", newDebt, debts, input)
	err = balanceRepo.SetDebt(ctx, tx, uid, newDebt)
	if err != nil {
		err = errors.New("failed update debts")
		Log.Error(err)
		return err
	}
	return nil
}
