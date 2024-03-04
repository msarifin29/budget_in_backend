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

type CreditController struct {
	CreditUsecase usecase.CreditUsecase
	Log           *logrus.Logger
}

func NewCreditController(CreditUsecase usecase.CreditUsecase, Log *logrus.Logger) *CreditController {
	return &CreditController{CreditUsecase: CreditUsecase, Log: Log}
}

func (c *CreditController) CreateCredit(ctx *gin.Context) {
	var req model.CreateCreditRequest
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
	res, err := c.CreditUsecase.CreateCredit(ctx, req)
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
		Data:    model.NewCredit(res),
	})
}

func (c *CreditController) UpdateCreditHistory(ctx *gin.Context) {
	var req model.UpdateHistoryCreditRequest
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
	params := model.UpdateHistoryCreditParams{
		Uid:         req.Uid,
		Id:          req.Id,
		CreditId:    req.CreditId,
		Status:      util.COMPLETED,
		TypePayment: req.TypePayment,
	}
	res, err := c.CreditUsecase.UpdateHistoryCredit(ctx, params)
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
