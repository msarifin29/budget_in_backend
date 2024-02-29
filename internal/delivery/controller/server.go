package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/msarifin29/be_budget_in/internal/config"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Log        *logrus.Logger
	Engine     *gin.Engine
	Con        config.Config
	TokenMaker util.Maker
	UserC      UserController
	ExpenseC   ExpenseController
}

func NewServer(Log *logrus.Logger, Con config.Config) (*Server, error) {
	db := config.Connection(Log)
	tokenMaker, err := util.NewJWTMaker(Con.TokenSymetricKey)

	if err != nil {
		Log.Fatalf("cannot generate token %t :", err)
	}

	// User
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo, Log, db)
	userController := NewUserController(userUsecase, Log, Con, tokenMaker)

	// Expense
	expenseRepo := repository.NewExpenseRepository()
	expenseUseCase := usecase.NewExpenseUsecase(expenseRepo, Log, db)
	expenseController := NewExpenseController(expenseUseCase, Log)

	server := &Server{
		Log:        Log,
		Con:        Con,
		TokenMaker: tokenMaker,
		UserC:      *userController,
		ExpenseC:   *expenseController,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", util.ValidCurrency)
		v.RegisterValidation("type_user", util.ValidType)
		v.RegisterValidation("expense_type", util.ValidExpenseType)
		v.RegisterValidation("category", util.ValidCategoryType)
		v.RegisterValidation("status", util.ValidStatusType)
	}

	server.setupRoute()
	return server, nil
}

func (server *Server) setupRoute() {

	router := gin.Default()

	binding.Validator.Engine()

	router.POST("/api/register", server.UserC.CreateUser)
	router.POST("/api/login", server.UserC.LoginUser)

	autRoutes := router.Group("/").Use(delivery.AuthMiddleware(server.TokenMaker))

	autRoutes.GET("/api/user/:uid", server.UserC.GetById)
	autRoutes.PUT("/api/update", server.UserC.UpdateUser)

	// Expense
	autRoutes.POST("api/expenses/create", server.ExpenseC.CreateExpense)
	autRoutes.GET("api/expenses/:id", server.ExpenseC.GetExpenseById)
	autRoutes.PUT("api/expenses/update", server.ExpenseC.UpdateExpense)
	autRoutes.DELETE("api/expenses/:id", server.ExpenseC.DeleteExpense)
	autRoutes.GET("api/expenses/", server.ExpenseC.GetExpenses)

	server.Engine = router
}

func (server *Server) Start(address string) error {
	return server.Engine.Run(address)
}
