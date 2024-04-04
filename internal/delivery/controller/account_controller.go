package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/sirupsen/logrus"
)

type AccountController struct {
	Usecase usecase.AccountUsacase
	Log     *logrus.Logger
}

func NewAccountController(Usecase usecase.AccountUsacase, Log *logrus.Logger) *AccountController {
	return &AccountController{Usecase: Usecase, Log: Log}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req model.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Errorf("failed binding request %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	res, err := c.Usecase.CreateAccount(ctx, req)
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
		Data:    res,
	})
}
func (c *AccountController) UpdateMaxBudget(ctx *gin.Context) {
	var req model.UpdateMaxBudgetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Errorf("failed binding request %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if req.Uid != authPayload.Uid {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	res, err := c.Usecase.UpdateMaxBudget(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: errors.New("failed update maximal budget").Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
func (c *AccountController) GetMaxBudget(ctx *gin.Context) {
	var req model.GetMaxBudgetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	if req.Uid != authPayload.Uid {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, model.MetaErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}
	param := model.GetMaxBudgetParam{
		Uid:       req.Uid,
		AccountId: req.AccountId,
	}
	res, err := c.Usecase.GetMaxBudget(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: errors.New("failed update maximal budget").Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
