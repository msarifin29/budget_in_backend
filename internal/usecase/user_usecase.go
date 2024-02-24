package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (model.UserResponse, error)
	GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error)
}

type UserUsecaseImpl struct {
	UserRepository repository.UserRepository
	Log            *logrus.Logger
	db             *sql.DB
}

func NewUserUsecase(UserRepository repository.UserRepository, Log *logrus.Logger, db *sql.DB) UserUsecase {
	return &UserUsecaseImpl{UserRepository: UserRepository, Log: Log, db: db}
}

func (u *UserUsecaseImpl) CreateUser(ctx context.Context, user model.CreateUserRequest) (model.UserResponse, error) {
	tx, err := u.db.Begin()
	defer tx.Rollback()

	var res model.UserResponse
	if err != nil {
		u.Log.Error(err)
		return res, err
	}

	password, hashErr := util.HashPassword(user.Password)

	if hashErr != nil {
		u.Log.Errorf("failed start transaction %e", hashErr)
		return res, hashErr
	}

	userReq := model.User{
		Uid:      uuid.New(),
		UserName: user.UserName,
		Email:    user.Email,
		Password: password,
		TypeUser: user.TypeUser,
		Balance:  user.Balance,
		Savings:  user.Savings,
		Cash:     user.Cash,
		Debts:    user.Debts,
		Currency: user.Currency,
	}

	req, reqErr := u.UserRepository.CreateUser(ctx, tx, userReq)
	if reqErr != nil {
		u.Log.Errorf("failed create user %e :", reqErr)
		return res, reqErr
	}

	cErr := tx.Commit()
	if cErr != nil {
		u.Log.Errorf("failed commit DB %e :", cErr)
		return res, cErr
	}

	return model.UserResponse{Uid: req.Uid, UserName: req.UserName}, nil
}

func (u *UserUsecaseImpl) GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error) {
	var res model.UserResponse
	tx, err := u.db.Begin()
	defer tx.Rollback()

	if err != nil {
		u.Log.Errorf("failed start transaction %e :", err)
		return res, err
	}
	ru, getErr := u.UserRepository.GetUser(ctx, tx, user.Email)
	fmt.Println("email => ", ru)

	if getErr != nil {
		u.Log.Errorf("user not found with email %e :", getErr)
		return res, getErr
	}

	passErr := util.CheckPassword(user.Password, ru.Password)
	if passErr != nil {
		u.Log.Errorf("invalid password %e :", passErr)
		return res, passErr
	}

	rum := model.UserResponse{Uid: ru.Uid, UserName: ru.UserName}

	return rum, nil
}
