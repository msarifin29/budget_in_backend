package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/msarifin29/be_budget_in/internal/config"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUsecase usecase.UserUsecase
	Log         *logrus.Logger
	Con         config.Config
	TokenMaker  util.Maker
}

func NewUserController(UserUsecase usecase.UserUsecase, Log *logrus.Logger, Con config.Config, TokenMaker util.Maker) *UserController {
	return &UserController{UserUsecase: UserUsecase, Log: Log, Con: Con, TokenMaker: TokenMaker}
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
		if strings.Contains(cErr.Error(), " Duplicate") {
			msg := fmt.Sprintf("The email address %s is already taken. Choose a unique email address to create your account.", userReq.Email)
			cErr = errors.New(msg)
			ctx.JSON(http.StatusUnprocessableEntity, model.MetaErrorResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: cErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: cErr.Error(),
		})
		return
	}

	token, _, ctErr := c.TokenMaker.CreateToken(res.UserName, c.Con.AccessTokenDuration, res.Uid)
	if ctErr != nil {
		c.Log.Error(ctErr)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
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

	token, _, ctErr := c.TokenMaker.CreateToken(res.UserName, c.Con.AccessTokenDuration, res.Uid)
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

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req model.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Errorf("binding %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	user, err := c.UserUsecase.GetById(ctx, req.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.MetaErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	err = c.UserUsecase.UpdateUser(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func (c *UserController) GetById(ctx *gin.Context) {
	var req model.UserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		c.Log.Errorf("binding %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	user, err := c.UserUsecase.GetById(ctx, req.Uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.MetaErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if user.UserName != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    user,
	})
}

func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var req model.EmailUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	ok, err := c.UserUsecase.ForgotPassword(ctx, req)
	if !ok || err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    ok,
	})
}

func (c *UserController) NonActivatedUser(ctx *gin.Context) {
	var req model.NonActiveUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Errorf("binding %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if authPayload.Uid != req.Uid {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}
	params := model.NonActiveUserParams{Uid: authPayload.Uid}
	ok, err := c.UserUsecase.NonActivatedUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    ok,
	})
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req model.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Error(err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if authPayload.Uid != req.Uid {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}
	ok, err := c.UserUsecase.ResetPassword(ctx, req)
	if !ok || err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    ok,
	})
}

func (c *UserController) CheckEmail(ctx *gin.Context) {
	var req model.CheckEmail

	if err := ctx.ShouldBindUri(&req); err != nil {
		c.Log.Errorf("binding %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	isValid, err := c.UserUsecase.GetEmailUser(ctx, req)
	if err != nil || !isValid {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    isValid,
	})
}
