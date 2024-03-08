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
		AccountId:   req.AccountId,
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

func (c *CreditController) GetAllCredit(ctx *gin.Context) {
	var req model.GetCreditsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	params := model.GetCreditParams{
		Uid:    authPayload.Uid,
		Limit:  req.TotalPage,
		Offset: (req.Page - 1) * req.TotalPage,
	}
	credits, total, err := c.CreditUsecase.GetAllCredit(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	lastPage := int32(total/float64(req.TotalPage) + 1)
	res := model.CreditsResponse{
		Page:      req.Page,
		TotalPage: req.TotalPage,
		LastPage:  lastPage,
		Total:     int32(total),
		Data:      credits,
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
func (c *CreditController) GetAllHistoriesCredit(ctx *gin.Context) {
	var req model.GetHistoriesCreditsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	params := model.GetHistoriesCreditParams{
		CreditId: req.CreditId,
		Limit:    req.TotalPage,
		Offset:   (req.Page - 1) * req.TotalPage,
	}
	histories, total, err := c.CreditUsecase.GetAllHistoryCredit(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	lastPage := int32(total/float64(req.TotalPage) + 1)
	res := model.HistoriesCreditsResponse{
		Page:      req.Page,
		TotalPage: req.TotalPage,
		LastPage:  lastPage,
		Total:     int32(total),
		Data:      histories,
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
