package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (model.UserResponse, error)
	GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error)
	UpdateUser(ctx context.Context, user model.UpdateUserRequest) error
	GetById(ctx context.Context, uid string) (model.AccountUser, error)
	ResetPassword(ctx context.Context, req model.EmailUserRequest) (bool, error)
	NonActivatedUser(ctx context.Context, req model.NonActiveUserParams) (bool, error)
}

type UserUsecaseImpl struct {
	UserRepository repository.UserRepository
	AccountRepo    repository.AccountRepository
	Log            *logrus.Logger
	db             *sql.DB
	conf           config.Config
}

func NewUserUsecase(UserRepository repository.UserRepository, AccountRepo repository.AccountRepository, Log *logrus.Logger, db *sql.DB, conf config.Config) UserUsecase {
	return &UserUsecaseImpl{UserRepository: UserRepository, AccountRepo: AccountRepo, Log: Log, db: db, conf: conf}
}

// NonActivatedUser implements UserUsecase.
func (u *UserUsecaseImpl) NonActivatedUser(ctx context.Context, req model.NonActiveUserParams) (bool, error) {
	tx, err := u.db.Begin()
	defer util.CommitOrRollback(tx)
	if err != nil {
		u.Log.Errorf("failed start transaction %e :", err)
		return false, err
	}
	user, err := u.UserRepository.GetUserAccount(ctx, tx, req.Uid)
	if err != nil {
		u.Log.Errorf("user not found with id %s :", err)
		message := "user not found with id :" + req.Uid
		return false, errors.New(message)
	}
	ok, err := u.UserRepository.NonActivatedUser(ctx, tx, user.Uid, util.NonActive)
	if !ok || err != nil {
		u.Log.Errorf("failed delete account %s :", err)
		err = errors.New("failed delete account")
		return false, err
	}
	return ok, nil
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
		u.Log.Errorf("failed hash password %e", hashErr)
		return res, hashErr
	}

	userReq := model.User{
		Uid:      uuid.New().String(),
		UserName: user.UserName,
		Email:    user.Email,
		Password: password,
		TypeUser: user.TypeUser,
		Balance:  0,
		Savings:  0,
		Cash:     0,
		Debts:    0,
		Currency: "IDR",
	}

	req, reqErr := u.UserRepository.CreateUser(ctx, tx, userReq)
	if reqErr != nil {
		u.Log.Errorf("failed create user %e :", reqErr)
		return res, reqErr
	}
	reqAccount := model.Account{
		UserId:      req.Uid,
		AccountId:   uuid.NewString(),
		AccountName: req.UserName,
		Balance:     user.Balance,
		Cash:        user.Cash,
		Debts:       0,
		Savings:     0,
		Currency:    "IDR",
	}
	account, errAc := u.AccountRepo.CreateAccount(ctx, tx, reqAccount)
	if errAc != nil {
		u.Log.Errorf("failed create account %e :", errAc)
		return res, errAc
	}
	return model.UserResponse{Uid: req.Uid, UserName: req.UserName, AccountId: account.AccountId}, nil
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
	account, errAcc := u.AccountRepo.GetAccountByUserId(ctx, tx, ru.Uid)
	if errAcc != nil {
		u.Log.Errorf("failed get account %e :", errAcc)
		return res, errAcc
	}
	rum := model.UserResponse{Uid: ru.Uid, UserName: ru.UserName, AccountId: account.AccountId}

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

func (u *UserUsecaseImpl) GetById(ctx context.Context, uid string) (model.AccountUser, error) {
	tx, _ := u.db.Begin()

	defer util.CommitOrRollback(tx)

	user, err := u.UserRepository.GetUserAccount(ctx, tx, uid)
	if err != nil {
		u.Log.Errorf("user not found with id %s :", err)
		message := "user not found with id :" + uid
		return model.AccountUser{}, errors.New(message)
	}
	return user, nil
}

func (u *UserUsecaseImpl) ResetPassword(ctx context.Context, req model.EmailUserRequest) (bool, error) {
	tx, _ := u.db.Begin()

	defer util.CommitOrRollback(tx)
	emailUser, username, err := u.UserRepository.GetUserByEmail(ctx, tx, req)
	if err != nil {
		err = fmt.Errorf("user not found with email %v", req.Email)
		u.Log.Error(err)
		return false, err
	}
	subject := "Reset Password"
	newPassword := subject + " : " + util.RandomString(6)

	receiver := emailUser
	r := util.NewRequest([]string{receiver}, subject, u.conf, u.Log)
	err = r.Send("../templates/email.html", map[string]string{"name": username, "password": util.RandomString(6)})
	if err != nil {
		return false, err
	}
	password, hashErr := util.HashPassword(newPassword)

	if hashErr != nil {
		u.Log.Errorf("failed hash password %e", hashErr)
		return false, hashErr
	}
	ok, errSetPass := u.UserRepository.UpdatePassword(ctx, tx, emailUser, password)
	if !ok || errSetPass != nil {
		errSetPass = errors.New("failed update password")
		return false, errSetPass
	}
	return true, nil
}
