package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (model.UserResponse, error)
	GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error)
	UpdateUser(ctx context.Context, user model.UpdateUserRequest) error
	GetById(ctx context.Context, uid string) (model.User, error)
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
	defer util.CommitOrRollback(tx)

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
		Uid:      uuid.New().String(),
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

	return model.UserResponse{Uid: req.Uid, UserName: req.UserName}, nil
}

func (u *UserUsecaseImpl) GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error) {
	var res model.UserResponse
	tx, err := u.db.Begin()
	defer util.CommitOrRollback(tx)

	if err != nil {
		u.Log.Errorf("failed start transaction %e :", err)
		return res, err
	}
	ru, getErr := u.UserRepository.GetUser(ctx, tx, user.Email)

	if getErr != nil {
		u.Log.Errorf("user not found with email %v :", user.Email)
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

func (u *UserUsecaseImpl) UpdateUser(ctx context.Context, user model.UpdateUserRequest) error {
	tx, _ := u.db.Begin()

	defer util.CommitOrRollback(tx)

	err := u.UserRepository.UpdateUserName(ctx, tx, user)

	if err != nil {
		u.Log.Errorf("failed update username %u :", err)
		return err
	}
	return nil
}

func (u *UserUsecaseImpl) GetById(ctx context.Context, uid string) (model.User, error) {
	tx, _ := u.db.Begin()

	defer util.CommitOrRollback(tx)

	user, err := u.UserRepository.GetById(ctx, tx, uid)
	if err != nil {
		u.Log.Errorf("user not found with id %s :", uid)
		message := "user not found with id :" + uid
		return model.User{}, errors.New(message)
	}
	return user, nil
}
