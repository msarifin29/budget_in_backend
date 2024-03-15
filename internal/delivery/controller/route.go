package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
)

func (server *Server) SetUpRoute() {

	router := gin.Default()

	binding.Validator.Engine()

	router.POST("/api/register", server.UserC.CreateUser)
	router.POST("/api/login", server.UserC.LoginUser)

	router.POST("/api/user/forgot_password", server.UserC.ForgotPassword)

	autRoutes := router.Group("/").Use(delivery.AuthMiddleware(server.TokenMaker))

	autRoutes.POST("/api/accounts/create", server.AccountC.CreateAccount)

	autRoutes.GET("/api/user/:uid", server.UserC.GetById)
	autRoutes.PUT("/api/update", server.UserC.UpdateUser)
	// Monthly reports
	autRoutes.GET("/api/user/monthly_report/", server.MonthReportC.GetMonthlyReport)

	// Expense
	autRoutes.POST("api/expenses/create", server.ExpenseC.CreateExpense)
	autRoutes.GET("api/expenses/:id", server.ExpenseC.GetExpenseById)
	autRoutes.PUT("api/expenses/update", server.ExpenseC.UpdateExpense)
	autRoutes.DELETE("api/expenses/:id", server.ExpenseC.DeleteExpense)
	autRoutes.GET("api/expenses/", server.ExpenseC.GetExpenses)
	autRoutes.GET("api/expenses/monthly_report/", server.ExpenseC.GetExpensesByMonth) //not used will be remove later

	// Incomes
	autRoutes.POST("/api/incomes/create", server.IncomeC.CreateIncome)
	autRoutes.GET("/api/incomes/", server.IncomeC.GetIncomes)
	autRoutes.GET("/api/incomes/monthly_report/", server.IncomeC.GetIncomesByMonth) //not used will be remove later

	// Credits
	autRoutes.POST("/api/credits/create", server.CreditC.CreateCredit)
	autRoutes.PUT("/api/credits/update_history", server.CreditC.UpdateCreditHistory)
	autRoutes.GET("/api/credits/", server.CreditC.GetAllCredit)
	autRoutes.GET("/api/histories_credits/", server.CreditC.GetAllHistoriesCredit)
	server.Engine = router
}
