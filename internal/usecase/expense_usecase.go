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

type ExpenseUsecase interface {
	CreateExpense(ctx context.Context, expense model.CreateExpenseParams) (model.Expense, error)
	GetExpenseById(ctx context.Context, request model.ExpenseParamWithId) (model.Expense, error)
	UpdateExpense(ctx context.Context, expense model.UpdateExpenseRequest) (bool, error)
	DeleteExpense(ctx context.Context, id float64) error
	GetExpenses(ctx context.Context, params model.GetExpenseParams) ([]model.ExpenseResponse, float64, error)
}

type ExpenseUsecaseImpl struct {
	CategoryRepo      repository.CategoryRepository
	ExpenseRepository repository.ExpenseRepository
	// BalanceRepository repository.BalanceRepository
	AccountRepo repository.AccountRepository
	Log         *logrus.Logger
	db          *sql.DB
}

// GetExpensez implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) GetExpenses(ctx context.Context, params model.GetExpenseParams) ([]model.ExpenseResponse, float64, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	total, err := u.ExpenseRepository.GetTotalExpenses(ctx, tx, params.Uid, params.Status, params.ExpenseType, params.Id, params.CreatedAt)
	if err != nil {
		u.Log.Errorf("failed get count expenses %v ", err)
		return []model.ExpenseResponse{}, 0, errors.New("failed get count expenses")
	}
	expenses, err := u.ExpenseRepository.GetExpenses(ctx, tx, params)
	if err != nil {
		u.Log.Errorf("failed get expenses %v ", err)
		return []model.ExpenseResponse{}, 0, errors.New("failed get expenses")
	}
	return expenses, total, nil
}

// CreateExpense implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) CreateExpense(ctx context.Context, expense model.CreateExpenseParams) (model.Expense, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)

	notes := zero.StringFromPtr(&expense.Notes)
	req := model.Expense{
		ExpenseType:   expense.ExpenseType,
		Total:         expense.Total,
		Status:        util.SUCCESS,
		Notes:         notes.String,
		Uid:           expense.Uid,
		TransactionId: uuid.NewString(),
		CreatedAt:     util.CreatedAt(expense.CreatedAt),
	}
	if req.Total < 2000 {
		u.Log.Errorf("Invalid input total min 2000 , actually %v", req.Total)
		inputErr := errors.New("invalid input total minimum")
		return model.Expense{}, inputErr
	}
	if expense.ExpenseType == util.DEBIT {
		err := util.NewBalance(ctx, tx, u.Log, u.AccountRepo, util.SUCCESS, expense.AccountId, req.Total)
		if err != nil {
			return model.Expense{}, err
		}
	} else if expense.ExpenseType == util.CASH {
		err := util.NewCash(ctx, tx, u.Log, u.AccountRepo, util.SUCCESS, expense.AccountId, req.Total)
		if err != nil {
			return model.Expense{}, err
		}
	}
	res, err := u.ExpenseRepository.CreateExpense(ctx, tx, req)
	if err != nil {
		u.Log.Errorf("failed create expense %e :", err)
		return model.Expense{}, err
	}
	paramCategory := model.Category{
		CategoryId: res.Id,
		Id:         expense.CategoryId,
		Title:      util.InputCategoryexpense(expense.CategoryId),
	}
	category, categoryErr := u.CategoryRepo.CreateCategoryExpense(ctx, tx, paramCategory)
	if categoryErr != nil {
		u.Log.Errorf("failed add category expense %e :", categoryErr)
		return model.Expense{}, categoryErr
	}
	update := zero.TimeFromPtr(res.UpdatedAt)
	return model.Expense{
		Uid:           res.Uid,
		Id:            res.Id,
		ExpenseType:   res.ExpenseType,
		Total:         res.Total,
		Category:      category.Title,
		Status:        res.Status,
		Notes:         notes.String,
		TransactionId: res.TransactionId,
		CreatedAt:     res.CreatedAt,
		UpdatedAt:     &update.Time,
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
	transaction := zero.StringFromPtr(&x.TransactionId)
	update := zero.TimeFromPtr(x.UpdatedAt)
	res := model.Expense{
		Uid:           x.Uid,
		Id:            x.Id,
		ExpenseType:   x.ExpenseType,
		Total:         x.Total,
		Category:      x.Category,
		Status:        x.Status,
		Notes:         notes.String,
		TransactionId: transaction.String,
		CreatedAt:     x.CreatedAt,
		UpdatedAt:     &update.Time,
	}
	return res, nil
}

// UpdateExpense implements ExpenseUsecase.
func (u *ExpenseUsecaseImpl) UpdateExpense(ctx context.Context, expense model.UpdateExpenseRequest) (bool, error) {
	tx, _ := u.db.Begin()
	defer util.CommitOrRollback(tx)
	x, err := u.ExpenseRepository.GetExpenseById(ctx, tx, expense.Id)
	if err != nil {
		u.Log.Errorf("expense not found with id %v :", expense.Id)
		return false, err
	}
	if x.Status != util.SUCCESS {
		return false, errors.New("status is cancelled")
	}

	if expense.ExpenseType == util.DEBIT && expense.ExpenseType == x.ExpenseType {
		err := util.NewBalance(ctx, tx, u.Log, u.AccountRepo, util.CANCELLED, expense.AccountId, x.Total)
		if err != nil {
			return false, err
		}
	} else if expense.ExpenseType == util.CASH && expense.ExpenseType == x.ExpenseType {
		err := util.NewCash(ctx, tx, u.Log, u.AccountRepo, util.CANCELLED, expense.AccountId, x.Total)
		if err != nil {
			return false, err
		}
	} else if expense.ExpenseType == "" || expense.ExpenseType != x.ExpenseType {
		return false, errors.New("invalid input type expense")
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
		return false, err
	}
	u.Log.Infof("success update expense with id %v", res.Id)
	return true, nil
}

func NewExpenseUsecase(
	CategoryRepo repository.CategoryRepository,
	ExpenseRepository repository.ExpenseRepository,
	// BalanceRepository repository.BalanceRepository,
	AccountRepo repository.AccountRepository,
	Log *logrus.Logger,
	db *sql.DB) ExpenseUsecase {
	return &ExpenseUsecaseImpl{
		CategoryRepo:      CategoryRepo,
		ExpenseRepository: ExpenseRepository,
		// BalanceRepository: BalanceRepository,
		AccountRepo: AccountRepo,
		Log:         Log,
		db:          db,
	}
}
