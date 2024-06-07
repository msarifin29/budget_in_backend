package usecase

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/msarifin29/be_budget_in/internal/model"
// 	"github.com/msarifin29/be_budget_in/internal/repository"
// 	"github.com/msarifin29/be_budget_in/util"
// 	"github.com/msarifin29/be_budget_in/util/zero"
// 	"github.com/sirupsen/logrus"
// )

// type CreditUsecase interface {
// 	CreateCredit(ctx context.Context, params model.CreateCreditRequest) (model.Credit, error)
// 	GetAllCredit(ctx context.Context, params model.GetCreditParams) ([]model.CreditResponse, float64, error)
// 	GetAllHistoryCredit(ctx context.Context, params model.GetHistoriesCreditParams) ([]model.HistoryCredit, float64, error)
// 	UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditParams) (model.UpdateHistoryResponse, error)
// }

// type CreditUsecaseImpl struct {
// 	CategoryRepo repository.CategoryRepository
// 	CreditRepo   repository.CreditRepository
// 	BalanceRepo  repository.BalanceRepository
// 	AccountRepo  repository.AccountRepository
// 	ExpenseRepo  repository.ExpenseRepository
// 	Log          *logrus.Logger
// 	db           *sql.DB
// }

// // GetAllHistoryCredit implements CreditUsecase.
// func (u *CreditUsecaseImpl) GetAllHistoryCredit(ctx context.Context, params model.GetHistoriesCreditParams) ([]model.HistoryCredit, float64, error) {
// 	tx, _ := u.db.Begin()
// 	defer util.CommitOrRollback(tx)

// 	historiesCredits, err := u.CreditRepo.GetAllHistoryCredit(ctx, tx, params)
// 	if err != nil {
// 		u.Log.Errorf("failed get all history credit %e", err)
// 		err = errors.New("failed get all history credit")
// 		return []model.HistoryCredit{}, 0, err
// 	}
// 	count, err := u.CreditRepo.GetCountHistoryCredit(ctx, tx, params.CreditId)
// 	if err != nil {
// 		u.Log.Errorf("failed count history credit %e", err)
// 		err = errors.New("failed get count history credit")
// 		return []model.HistoryCredit{}, 0, err
// 	}
// 	return historiesCredits, count, nil
// }

// // CreateCredit implements CreditUsecase.
// func (u *CreditUsecaseImpl) CreateCredit(ctx context.Context, params model.CreateCreditRequest) (model.Credit, error) {
// 	loanTerm := util.GetTotalMonth(util.Date(params.StartDate), util.Date(params.EndDate)) + 1
// 	tx, _ := u.db.Begin()
// 	defer util.CommitOrRollback(tx)
// 	req := model.Credit{
// 		Uid:            params.Uid,
// 		CategoryCredit: params.CategoryCredit,
// 		TypeCredit:     params.TypeCredit,
// 		Total:          float64(loanTerm) * params.Installment,
// 		LoanTerm:       float64(loanTerm),
// 		StatusCredit:   util.ACTIVE,
// 		Installment:    params.Installment,
// 		PaymentTime:    util.Date(params.StartDate).Day(),
// 		StartDate:      util.Date(params.StartDate),
// 		EndDate:        util.Date(params.EndDate),
// 	}
// 	creditRes, err := u.CreditRepo.CreateCredit(ctx, tx, req)
// 	if err != nil {
// 		u.Log.Errorf("failed create credit %v", err)
// 		return model.Credit{}, err
// 	}

// 	err = NewHistoryCredit(ctx, tx, u.CreditRepo, creditRes, u.Log)
// 	if err != nil {
// 		u.Log.Errorf("failed create history credit %v", err)
// 		return model.Credit{}, err
// 	}

// 	err = NewDebts(ctx, tx, u.Log, u.BalanceRepo, util.ACTIVE, creditRes.Uid, creditRes.Total)
// 	if err != nil {
// 		u.Log.Errorf("failed update depts %v", err)
// 		return model.Credit{}, err
// 	}
// 	paramCategory := model.Category{
// 		CategoryId: creditRes.Id,
// 		Id:         float64(params.CategoryId),
// 		Title:      util.InputCategoryCredit(float64(params.CategoryId)),
// 	}
// 	category, catErr := u.CategoryRepo.CreateCategoryCredits(ctx, tx, paramCategory)
// 	if catErr != nil {
// 		u.Log.Errorf("failed create category credit %v", catErr)
// 		return model.Credit{}, catErr
// 	}
// 	update := zero.TimeFromPtr(creditRes.UpdatedAt)
// 	res := model.Credit{
// 		Uid:            creditRes.Uid,
// 		Id:             creditRes.Id,
// 		CategoryCredit: category.Title,
// 		TypeCredit:     creditRes.TypeCredit,
// 		Total:          creditRes.Total,
// 		LoanTerm:       float64(loanTerm),
// 		StatusCredit:   util.ACTIVE,
// 		Installment:    creditRes.Installment,
// 		CreatedAt:      creditRes.CreatedAt,
// 		UpdatedAt:      &update.Time,
// 		PaymentTime:    creditRes.PaymentTime,
// 	}
// 	return res, nil
// }

// // UpdateCredit implements CreditUsecase.
// func (u *CreditUsecaseImpl) GetAllCredit(ctx context.Context, params model.GetCreditParams) ([]model.CreditResponse, float64, error) {
// 	tx, _ := u.db.Begin()
// 	defer util.CommitOrRollback(tx)

// 	credits, err := u.CreditRepo.GetAllCredit(ctx, tx, params)
// 	if err != nil {
// 		u.Log.Errorf("failed get all credit %v", err)
// 		err = errors.New("failed get all credit")
// 		return []model.CreditResponse{}, 0, err
// 	}
// 	count, err := u.CreditRepo.GetCountCredit(ctx, tx, params.Uid)
// 	if err != nil {
// 		u.Log.Errorf("failed get count credit %v", err)
// 		err = errors.New("failed get count credit")
// 		return []model.CreditResponse{}, 0, err
// 	}
// 	return credits, count, nil
// }

// // UpdateHistoryCredit implements CreditUsecase.
// func (u *CreditUsecaseImpl) UpdateHistoryCredit(ctx context.Context, params model.UpdateHistoryCreditParams) (model.UpdateHistoryResponse, error) {
// 	tx, _ := u.db.Begin()
// 	defer util.CommitOrRollback(tx)

// 	credit, err := u.CreditRepo.GetCreditById(ctx, tx, model.GetCreditRequest{Uid: params.Uid, Id: params.CreditId})
// 	if err != nil {
// 		u.Log.Errorf("failed get credit with credit id %v", params.CreditId)
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	if params.Status == util.ACTIVE {
// 		err := errors.New("cannot update history credit with same value")
// 		return model.UpdateHistoryResponse{}, err
// 	}

// 	reqId := model.GetHistoryCreditRequest{Uid: params.Uid, Id: params.Id}
// 	historyC, err := u.CreditRepo.GetHistoryCreditById(ctx, tx, reqId)
// 	if err != nil {
// 		u.Log.Errorf("failed get history credit with id %v", params.Id)
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	if historyC.Status != util.ACTIVE {
// 		err = errors.New("status credit is completed")
// 		u.Log.Error(err)
// 		return model.UpdateHistoryResponse{}, err
// 	}

// 	ok, err := u.CreditRepo.UpdateHistoryCredit(ctx, tx, params)
// 	if err != nil || !ok {
// 		u.Log.Error(err)
// 		err = errors.New("failed update history credit")
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	err = UpdateTotalBalanceOrCash(ctx, tx, u.CreditRepo, u.BalanceRepo, historyC, u.AccountRepo, params.TypePayment, params.AccountId, u.Log)
// 	if err != nil {
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	newCredit := credit.Total - historyC.Total
// 	ok, totalErr := u.CreditRepo.UpdateTotalCredit(ctx, tx, params.Uid, params.CreditId, newCredit)
// 	if totalErr != nil || !ok {
// 		u.Log.Error(err)
// 		err = errors.New("failed update total credit")
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	isCompleted, errCheck := checkHistoryCredit(ctx, tx, u.CreditRepo, params.CreditId)
// 	if errCheck != nil {
// 		u.Log.Error(errCheck)
// 		errCheck = errors.New("failed get history credit")
// 		return model.UpdateHistoryResponse{}, errCheck
// 	}
// 	if isCompleted {
// 		ok, errUp := u.CreditRepo.UpdateCredit(ctx, tx, model.UpdateCreditRequest{Uid: params.Uid, Id: params.CreditId, StatusCredit: util.COMPLETED})
// 		if !ok || err != nil {
// 			u.Log.Error(errUp)
// 			errUp = errors.New("failed update credit")
// 			return model.UpdateHistoryResponse{}, errUp
// 		}
// 	}

// 	// Update debts from user
// 	newCreditCompletted := credit.Installment * credit.LoanTerm
// 	err = util.NewDebts(ctx, tx, u.Log, u.AccountRepo, util.COMPLETED, params.AccountId, newCreditCompletted)
// 	if err != nil {
// 		u.Log.Error(err)
// 		err = errors.New("failed update debts user")
// 		return model.UpdateHistoryResponse{}, err
// 	}
// 	now := time.Now()
// 	expense := model.Expense{
// 		ExpenseType:   params.TypePayment,
// 		Total:         historyC.Total,
// 		Category:      util.COSTANDBILL,
// 		Status:        util.SUCCESS,
// 		Uid:           params.Uid,
// 		TransactionId: uuid.NewString(),
// 		CreatedAt:     &now,
// 	}
// 	newEx, er := u.ExpenseRepo.CreateExpense(ctx, tx, expense)
// 	if er != nil {
// 		u.Log.Error(er)
// 		er = errors.New("failed create expense")
// 		return model.UpdateHistoryResponse{}, er
// 	}
// 	paramCategory := model.Category{
// 		CategoryId: newEx.Id,
// 		Id:         8,
// 		Title:      util.COSTANDBILL,
// 	}
// 	_, categoryErr := u.CategoryRepo.CreateCategoryExpense(ctx, tx, paramCategory)
// 	if categoryErr != nil {
// 		u.Log.Errorf("failed add category expense %e :", categoryErr)
// 		return model.UpdateHistoryResponse{}, categoryErr
// 	}
// 	return model.UpdateHistoryResponse{
// 		Id:          historyC.Id,
// 		Th:          historyC.Th,
// 		Total:       historyC.Total,
// 		Status:      params.Status,
// 		TypePayment: params.TypePayment,
// 		CreatedAt:   historyC.CreatedAt,
// 	}, nil
// }

// func NewCreditUsecase(CreditRepo repository.CreditRepository, BalanceRepo repository.BalanceRepository,
// 	AccountRepo repository.AccountRepository, Log *logrus.Logger, db *sql.DB, ExpenseRepo repository.ExpenseRepository,
// 	CategoryRepo repository.CategoryRepository) CreditUsecase {
// 	return &CreditUsecaseImpl{CreditRepo: CreditRepo, BalanceRepo: BalanceRepo, AccountRepo: AccountRepo, Log: Log, db: db,
// 		ExpenseRepo: ExpenseRepo, CategoryRepo: CategoryRepo}
// }

// func NewHistoryCredit(ctx context.Context, tx *sql.Tx, creditRepo repository.CreditRepository, credit model.Credit, log *logrus.Logger) error {
// 	dates := util.GenerateDates(*credit.StartDate, *credit.EndDate)
// 	for i := 0; i < len(dates); i++ {
// 		date := util.Date(dates[i])
// 		log.Info(date)
// 		req := model.HistoryCredit{
// 			CreditId:    credit.Id,
// 			Th:          float64(i + 1),
// 			Total:       credit.Installment,
// 			Status:      util.ACTIVE,
// 			TypePayment: "",
// 			PaymentTime: credit.StartDate.Day(),
// 			Date:        date,
// 		}
// 		hc, err := creditRepo.CreateHistoryCredit(ctx, tx, req)
// 		log.Printf("history credit %v ", hc)
// 		if err != nil {
// 			log.Errorf("history error %v ", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func UpdateTotalBalanceOrCash(ctx context.Context, tx *sql.Tx,
// 	creditRepo repository.CreditRepository,
// 	BalanceRepo repository.BalanceRepository,
// 	historyCredit model.HistoryCredit,
// 	accountRepo repository.AccountRepository,
// 	typePayment string,
// 	accountId string,
// 	Log *logrus.Logger,
// ) error {
// 	switch typePayment {
// 	case util.DEBIT:
// 		account, err := accountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
// 		balance := account.Balance
// 		if err != nil {
// 			err = fmt.Errorf("failed get balance from user id %v", account.UserId)
// 			return err
// 		}
// 		if balance < historyCredit.Total {
// 			err = errors.New("cannot upgrade balance with total greater than balance")
// 			return err
// 		}
// 		newBalance := balance - historyCredit.Total
// 		Log.Infof("newbalance = %v, balance = %v, input = %v", newBalance, balance, historyCredit.Total)
// 		err = accountRepo.UpdateAccountBalance(ctx, tx, model.UpdateAccountBalance{AccountId: accountId, Balance: newBalance})
// 		if err != nil {
// 			err = errors.New("failed update balance")
// 			return err
// 		}
// 	case util.CASH:
// 		account, err := accountRepo.GetAccountByAccountId(ctx, tx, model.GetAccountRequest{AccountId: accountId})
// 		cash := account.Cash
// 		if err != nil {
// 			err = fmt.Errorf("failed get cash from user id %v", account.UserId)
// 			return err
// 		}
// 		if cash < historyCredit.Total {
// 			err = errors.New("cannot upgrade cash with total greater than cash")
// 			return err
// 		}
// 		newCash := cash - historyCredit.Total
// 		Log.Infof("newCash = %v, cash = %v, input = %v", newCash, cash, historyCredit.Total)
// 		err = accountRepo.UpdateAccountCash(ctx, tx, model.UpdateAccountCash{AccountId: accountId, Cash: newCash})
// 		if err != nil {
// 			err = errors.New("failed update cash")
// 			return err
// 		}
// 	}
// 	return nil
// }

// func checkHistoryCredit(ctx context.Context, tx *sql.Tx, creditRepo repository.CreditRepository, creditId float64) (bool, error) {
// 	historyCredits, err := creditRepo.GetAllHistoryCreditById(ctx, tx, creditId)
// 	if err != nil {
// 		err = fmt.Errorf("failed get history credit %v", creditId)
// 		return false, err
// 	}
// 	for _, i := range historyCredits {
// 		fmt.Println("status credit => ", i.Status)
// 		if i.Status == util.ACTIVE {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }
