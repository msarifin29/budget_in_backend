package test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewUserRepository()
	usecase := usecase.NewUserUsecase(repo, log, db)
	user := model.CreateUserRequest{
		UserName: "hash",
		Email:    "has@mail.com",
		Password: "password",
		TypeUser: "personal",
	}
	_, err := usecase.CreateUser(context.Background(), user)
	assert.NoError(t, err)
}
func TestGetUser(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewUserRepository()
	usecase := usecase.NewUserUsecase(repo, log, db)
	user := model.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "123456",
	}
	u, err := usecase.GetUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
}
func TestUpdateUser(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewUserRepository()
	usecase := usecase.NewUserUsecase(repo, log, db)
	user := model.UpdateUserRequest{
		UserName: "testing",
		Uid:      "f1687230-49d3-4657-96be-9b934ed0387f",
	}
	err := usecase.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestGetUserById(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewUserRepository()
	usecase := usecase.NewUserUsecase(repo, log, db)
	uid := "f1687230-49d3-4657-96be-9b934ed0387f"

	u, err := usecase.GetById(context.Background(), uid)
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
}
func TestCreateExpense(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo, log, db)
	x := model.CreateExpenseRequest{
		Uid:         "deb3823d-5581-4e98-896c-06e5aa3bac4a",
		ExpenseType: "Credit",
		Total:       1000000,
	}
	res, err := usecase.CreateExpense(context.Background(), x)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
func TestGetExpense(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo, log, db)
	res, err := usecase.GetExpenseById(context.Background(), model.ExpenseParamWithId{Id: 9})
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
func TestUpdateExpense(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo, log, db)
	x := model.UpdateExpenseRequest{Id: 4, Total: 20500}
	res, err := usecase.UpdateExpense(context.Background(), x)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
func TestDeleteExpense(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo, log, db)
	err := usecase.DeleteExpense(context.Background(), 4)
	assert.NoError(t, err)
}
func TestGetExpenses(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo, log, db)

	params := model.GetExpenseParams{
		// Uid:       "f1687230-49d3-4657-96be-9b934ed0387f",
		Limit:  1,
		Offset: 2,
	}
	x, err := usecase.GetExpenses(context.Background(), params)
	assert.NoError(t, err)
	assert.NotEmpty(t, x)
}
