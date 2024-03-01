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

type ExpenseUsecase interface {
	CreateExpense(ctx context.Context, expense model.CreateExpenseRequest) (model.Expense, error)
	GetExpenseById(ctx context.Context, request model.ExpenseParamWithId) (model.Expense, error)
	UpdateExpense(ctx context.Context, expense model.UpdateExpenseRequest) (model.Expense, error)
	DeleteExpense(ctx context.Context, id float64) error
	GetExpenses(ctx context.Context, params model.GetExpenseParams) ([]model.Expense, float64, error)
}

type ExpenseUsecaseImpl struct {
	ExpenseRepository repository.ExpenseRepository
	BalanceRepository repository.BalanceRepository
	Log               *logrus.Logger
	db                *sql.DB
}

// GetExpensez implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) GetExpenses(ctx context.Context, params model.GetExpenseParams) ([]model.Expense, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	total, err := u.ExpenseRepository.GetTotalExpenses(ctx, tx, params.Uid, params.Status)
	if err != nil {
		u.Log.Errorf("failed get toal expenses %v ", err)
		return []model.Expense{}, 0, err
	}
	expenses, err := u.ExpenseRepository.GetExpenses(ctx, tx, params)
	if err != nil {
		u.Log.Errorf("failed get expenses %v ", err)
		return []model.Expense{}, 0, err
	}
	return expenses, total, nil
}

// CreateExpense implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) CreateExpense(ctx context.Context, expense model.CreateExpenseRequest) (model.Expense, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	notes := zero.StringFromPtr(&expense.Notes)
	req := model.Expense{
		ExpenseType: expense.ExpenseType,
		Total:       expense.Total,
		Category:    expense.Category,
		Status:      util.SUCCESS,
		Notes:       notes.String,
		Uid:         expense.Uid,
	}
	if req.Total < 2000 {
		u.Log.Errorf("Invalid input total min 2000 , actually %v", req.Total)
		inputErr := errors.New("invalid input total minimum")
		return model.Expense{}, inputErr
	}
	if expense.ExpenseType == util.DEBIT {
		err := NewBalance(ctx, tx, u.Log, u.BalanceRepository, util.SUCCESS, expense.Uid, req.Total)
		if err != nil {
			return model.Expense{}, err
		}
	} else if expense.ExpenseType == util.CASH {
		err := NewCash(ctx, tx, u.Log, u.BalanceRepository, util.SUCCESS, expense.Uid, req.Total)
		if err != nil {
			return model.Expense{}, err
		}
	} else if expense.ExpenseType == util.CREDIT {
		err := NewDebts(ctx, tx, u.Log, u.BalanceRepository, expense.Uid, req.Total)
		if err != nil {
			return model.Expense{}, err
		}
	}
	res, err := u.ExpenseRepository.CreateExpense(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create expense %e :", err)
		return model.Expense{}, err
	}

	update := zero.TimeFrom(res.UpdatedAt)
	return model.Expense{
		Uid:         res.Uid,
		Id:          res.Id,
		ExpenseType: res.ExpenseType,
		Total:       res.Total,
		Category:    res.Category,
		Status:      res.Status,
		Notes:       notes.String,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   update.Time,
	}, nil
}

// DeleteExpense implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) DeleteExpense(ctx context.Context, id float64) error {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	err := u.ExpenseRepository.DeleteExpense(ctx, tx, id)
	return err
}

// GetExpenseById implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) GetExpenseById(ctx context.Context, request model.ExpenseParamWithId) (model.Expense, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	x, err := u.ExpenseRepository.GetExpenseById(ctx, tx, request.Id)
	if err != nil {
		u.Log.Errorf("expense not found with id %v :", request.Id)
		return model.Expense{}, err
	}
	notes := zero.StringFromPtr(&x.Notes)
	update := zero.TimeFromPtr(&x.UpdatedAt)
	res := model.Expense{
		Uid:         x.Uid,
		Id:          x.Id,
		ExpenseType: x.ExpenseType,
		Total:       x.Total,
		Category:    x.Category,
		Status:      x.Status,
		Notes:       notes.String,
		CreatedAt:   x.CreatedAt,
		UpdatedAt:   update.Time,
	}
	return res, nil
}

// UpdateExpense implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) UpdateExpense(ctx context.Context, expense model.UpdateExpenseRequest) (model.Expense, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	x, err := u.ExpenseRepository.GetExpenseById(ctx, tx, expense.Id)
	if err != nil {
		u.Log.Errorf("expense not found with id %v :", expense.Id)
		return model.Expense{}, err
	}

	if expense.ExpenseType == util.DEBIT {
		err := NewBalance(ctx, tx, u.Log, u.BalanceRepository, util.CANCELLED, x.Uid, x.Total)
		if err != nil {
			return model.Expense{}, err
		}
	} else if expense.ExpenseType == util.CASH {
		err := NewCash(ctx, tx, u.Log, u.BalanceRepository, util.CANCELLED, x.Uid, x.Total)
		if err != nil {
			return model.Expense{}, err
		}
	} else if expense.ExpenseType == util.CREDIT {
		return model.Expense{}, errors.New("invalid input type expense")
	}

	x.Status = util.CANCELLED

	req := model.Expense{
		Id:          x.Id,
		Status:      x.Status,
		ExpenseType: expense.ExpenseType,
	}

	res, err := u.ExpenseRepository.UpdateExpense(ctx, tx, req, expense.Id)
	if err != nil {
		u.Log.Errorf("failed update expense %u :", err)
		return model.Expense{}, err
	}
	return model.Expense{
		Uid:         res.Uid,
		Id:          res.Id,
		ExpenseType: res.ExpenseType,
		Total:       res.Total,
		Notes:       res.Notes,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func NewExpenseUsecase(
	ExpenseRepository repository.ExpenseRepository,
	BalanceRepository repository.BalanceRepository,
	Log *logrus.Logger,
	db *sql.DB) ExpenseUsecase {
	return &ExpenseUsecaseImpl{
		ExpenseRepository: ExpenseRepository,
		BalanceRepository: BalanceRepository,
		Log:               Log,
		db:                db,
	}
}
