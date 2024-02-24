package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUsecase usecase.UserUsecase
	Log         *logrus.Logger
	Con         config.Config
}

func NewUserController(UserUsecase usecase.UserUsecase, Log *logrus.Logger, Con config.Config) *UserController {
	return &UserController{UserUsecase: UserUsecase, Log: Log, Con: Con}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var userReq model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		c.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	res, cErr := c.UserUsecase.CreateUser(ctx, userReq)

	if cErr != nil {
		e := errors.New("UNIQUE constraint failed")
		if errors.As(cErr, &e) {
			ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
				Code:    http.StatusBadRequest,
				Message: cErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: cErr.Error(),
		})
		return
	}

	// Token
	tokenMaker, tokenErr := util.NewJWTMaker(c.Con.TokenSymetricKey)
	if tokenErr != nil {
		c.Log.Error(tokenErr)
		ctx.JSON(http.StatusInternalServerError, model.MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: errors.New("cannot generate token").Error(),
		})
		return
	}

	token, _, ctErr := tokenMaker.CreateToken(res.UserName, c.Con.AccessTokenDuration)
	if ctErr != nil {
		c.Log.Error(ctErr)
		ctx.JSON(http.StatusInternalServerError, model.MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: errors.New("cannot generate token").Error(),
		})
		return
	}
	tokenUser := model.TokenUserResponse{Token: token, UserRes: res}

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    tokenUser,
	})
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var userReq model.LoginUserRequest
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		c.Log.Errorf("binding %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, getErr := c.UserUsecase.GetUser(ctx, userReq)
	if getErr != nil {
		ctx.JSON(http.StatusNotFound, model.MetaErrorResponse{
			Code:    http.StatusNotFound,
			Message: errors.New("invalid email or password").Error(),
		})
		return
	}

	tokenMaker, tokenErr := util.NewJWTMaker(c.Con.TokenSymetricKey)
	if tokenErr != nil {
		c.Log.Errorf("new token %t:", tokenErr)
		ctx.JSON(http.StatusInternalServerError, model.MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: errors.New("cannot generate token").Error(),
		})
		return
	}

	token, _, ctErr := tokenMaker.CreateToken(res.UserName, c.Con.AccessTokenDuration)
	if ctErr != nil {
		c.Log.Errorf("cannot generate token %t:", ctErr)
		ctx.JSON(http.StatusInternalServerError, model.MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: errors.New("cannot generate token").Error(),
		})
		return
	}
	tokenUser := model.TokenUserResponse{Token: token, UserRes: res}

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    tokenUser,
	})
}
