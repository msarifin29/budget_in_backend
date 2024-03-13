package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/google/uuid"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

const (
	CONFIG_SMTP_HOST = "smtp.gmail.com"
	CONFIG_SMTP_PORT = 587
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user model.CreateUserRequest) (model.UserResponse, error)
	GetUser(ctx context.Context, user model.LoginUserRequest) (model.UserResponse, error)
	UpdateUser(ctx context.Context, user model.UpdateUserRequest) error
	GetById(ctx context.Context, uid string) (model.User, error)
	ResetPassword(ctx context.Context, req model.EmailUserRequest) (bool, error)
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
		AccountName: "",
		Balance:     user.Balance,
		Cash:        user.Cash,
		Debts:       0,
		Savings:     0,
		Currency:    "IDR",
	}
	_, err = u.AccountRepo.CreateAccount(ctx, tx, reqAccount)
	if err != nil {
		u.Log.Errorf("failed create account %e :", err)
		return res, err
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

	user, err := u.UserRepository.GetUserAccount(ctx, tx, uid)
	if err != nil {
		u.Log.Errorf("user not found with id %s :", err)
		message := "user not found with id :" + uid
		return model.User{}, errors.New(message)
	}
	return user, nil
}

func (u *UserUsecaseImpl) ResetPassword(ctx context.Context, req model.EmailUserRequest) (bool, error) {
	tx, _ := u.db.Begin()

	defer util.CommitOrRollback(tx)
	emailUser, err := u.UserRepository.GetUserByEmail(ctx, tx, req)
	if err != nil {
		err = fmt.Errorf("user not found with email %v", req.Email)
		u.Log.Error(err)
		return false, err
	}
	to := []string{emailUser}
	cc := []string{}
	subject := "New Password"
	newPassword := subject + " : " + util.RandomString(6)
	err = sendMail(to, cc, subject, newPassword, u.conf)
	if err != nil {
		u.Log.Errorf("failed send email %v:", err)
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

func sendMail(to []string, cc []string, subject, message string, conf config.Config) error {
	body := "From: " + conf.SenderName + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", conf.AuthEmail, conf.AuthPassword, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, conf.AuthEmail, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
