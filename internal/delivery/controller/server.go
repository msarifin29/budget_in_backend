package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Log    *logrus.Logger
	Engine *gin.Engine
	Con    config.Config
	UserC  UserController
}

func NewServer(Log *logrus.Logger, Con config.Config) (*Server, error) {
	db := config.Connection(Log)

	// User
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo, Log, db)
	userController := NewUserController(userUsecase, Log, Con)

	server := &Server{
		Log:   Log,
		Con:   Con,
		UserC: *userController,
	}

	server.setupRoute()
	return server, nil
}

func (server *Server) setupRoute() {
	router := gin.Default()

	binding.Validator.Engine()

	router.POST("/api/register", server.UserC.CreateUser)
	router.POST("/api/login", server.UserC.LoginUser)

	server.Engine = router
}

func (server *Server) Start(address string) error {
	return server.Engine.Run(address)
}
