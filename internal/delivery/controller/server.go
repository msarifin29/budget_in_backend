package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Log          *logrus.Logger
	Engine       *gin.Engine
	Con          config.Config
	TokenMaker   util.Maker
	UserC        UserController
	ExpenseC     ExpenseController
	IncomeC      IncomeController
	AccountC     AccountController
	MonthReportC MonthlyReportController
	PrivacyC     PrivacyController
}

func NewServer(Log *logrus.Logger, Con config.Config) (*Server, error) {
	db := config.Connection(Log)
	tokenMaker, err := util.NewJWTMaker(Con.TokenSymetricKey)

	if err != nil {
		Log.Fatalf("cannot generate token %t :", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository()
	expenseRepo := repository.NewExpenseRepository()
	balanceRepo := repository.NewBalanceRepository()
	incomeRepo := repository.NewIncomeRepository()
	accountRepo := repository.NewAccountRepository()
	monthlyRepo := repository.NewMonthlyRepository()
	categoryRepo := repository.NewCategoryRepository()

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, accountRepo, Log, db, Con)
	expenseUseCase := usecase.NewExpenseUsecase(categoryRepo, expenseRepo, balanceRepo, accountRepo, Log, db)
	incomeUsecase := usecase.NewIncomeUsecase(incomeRepo, balanceRepo, accountRepo, Log, db, categoryRepo)
	accountUsacase := usecase.NewAccountUsacase(accountRepo, expenseRepo, Log, db)
	monthlyUsecase := usecase.NewMonthlyReportUsecase(monthlyRepo, Log, db)

	// Controller
	userController := NewUserController(userUsecase, Log, Con, tokenMaker)
	expenseController := NewExpenseController(expenseUseCase, Log)
	incomeController := NewIncomeController(incomeUsecase, Log)
	accountController := NewAccountController(accountUsacase, Log)
	monthlyController := NewMonthlyController(monthlyUsecase, Log)
	privacyController := NewPrivacyController(Log)

	server := &Server{
		Log:          Log,
		Con:          Con,
		TokenMaker:   tokenMaker,
		UserC:        *userController,
		ExpenseC:     *expenseController,
		IncomeC:      *incomeController,
		AccountC:     *accountController,
		MonthReportC: *monthlyController,
		PrivacyC:     *privacyController,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", util.ValidCurrency)
		v.RegisterValidation("type_user", util.ValidType)
		v.RegisterValidation("expense_type", util.ValidExpenseType)
		v.RegisterValidation("category", util.ValidCategoryType)
		v.RegisterValidation("status", util.ValidStatusType)
		v.RegisterValidation("type_income", util.ValidIncomeType)
		v.RegisterValidation("category_income", util.ValidCategoryType)
		v.RegisterValidation("type_credit", util.ValidTypeCredit)
		v.RegisterValidation("status_credit", util.ValidStatusHistoryCredit)
		v.RegisterValidation("category_credit", util.ValidCategoryCredit)
		v.RegisterValidation("category_expense", util.ValidCategoryExpense)
	}

	server.SetUpRoute()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Engine.Run(address)
}
