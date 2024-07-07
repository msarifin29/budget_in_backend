package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
)

func (server *Server) SetUpRoute() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	binding.Validator.Engine()

	// No Authorization
	router.POST("/api/register", server.UserC.CreateUser)
	router.POST("/api/login", server.UserC.LoginUser)
	router.GET("/api/check-email/:email", server.UserC.CheckEmail)
	router.POST("/api/user/forgot_password", server.UserC.ForgotPassword)
	router.GET("/api/privacy-police/:lang", server.PrivacyC.PrivacyPolice)
	router.PUT("/api/user/delete/:email", server.UserC.DeleteEmailUserUser)

	autRoutes := router.Group("/").Use(delivery.AuthMiddleware(server.TokenMaker))

	// Accounts
	autRoutes.POST("/api/accounts/create", server.AccountC.CreateAccount)
	autRoutes.GET("/api/accounts/", server.AccountC.GetAllAccounts)
	autRoutes.PUT("/api/accounts/update_max_budget", server.AccountC.UpdateMaxBudget)
	autRoutes.GET("/api/accounts/max_budget/", server.AccountC.GetMaxBudget)

	// Users
	autRoutes.GET("/api/user/:uid", server.UserC.GetById)
	autRoutes.PUT("/api/update", server.UserC.UpdateUser)
	autRoutes.PUT("/api/user/delete", server.UserC.NonActivatedUser)
	autRoutes.PUT("/api/user/reset_password", server.UserC.ResetPassword)

	// Monthly reports
	autRoutes.GET("/api/user/monthly_report/:uid", server.MonthReportC.GetMonthlyReport)
	autRoutes.GET("/api/user/monthly-report-detail/", server.MonthReportC.GetMonthlyReportDetail)
	autRoutes.GET("/api/user/monthly-report/category/", server.MonthReportC.GetMonthlyReportCategory)

	// Expense
	autRoutes.POST("api/expenses/create", server.ExpenseC.CreateExpense)
	autRoutes.GET("api/expenses/:id", server.ExpenseC.GetExpenseById)
	autRoutes.PUT("api/expenses/update", server.ExpenseC.UpdateExpense)
	autRoutes.DELETE("api/expenses/:id", server.ExpenseC.DeleteExpense)
	autRoutes.GET("api/expenses/", server.ExpenseC.GetExpenses)

	// Incomes
	autRoutes.POST("/api/incomes/create", server.IncomeC.CreateIncome)
	autRoutes.GET("/api/incomes/", server.IncomeC.GetIncomes)
	autRoutes.PUT("/api/incomes/cash-withdrawal", server.IncomeC.CashWithdrawal)

	server.Engine = router
}
