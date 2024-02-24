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
