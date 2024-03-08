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

type IncomeController struct {
	IncomeUsecase usecase.IncomeUsecase
	Log           *logrus.Logger
}

func NewIncomeController(IncomeUsecase usecase.IncomeUsecase, Log *logrus.Logger) *IncomeController {
	return &IncomeController{IncomeUsecase: IncomeUsecase, Log: Log}
}

func (c *IncomeController) CreateIncome(ctx *gin.Context) {
	var req model.CreateIncomeRequest
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
	res, err := c.IncomeUsecase.CreateIncome(ctx, req)
	if err != nil {
		c.Log.Errorf("failed create income %t:", err)
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
func (c *IncomeController) GetIncomes(ctx *gin.Context) {
	var req model.GetIncomeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	params := model.GetIncomeParams{
		Uid:            authPayload.Uid,
		CategoryIncome: req.CategoryIncome,
		Limit:          req.TotalPage,
		Offset:         (req.Page - 1) * req.TotalPage,
	}
	incomes, total, err := c.IncomeUsecase.GetIncomes(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	lastPage := int32(total/float64(req.TotalPage) + 1)
	res := model.IncomesResponse{
		Page:      req.Page,
		TotalPage: req.TotalPage,
		LastPage:  lastPage,
		Total:     int32(total),
		Data:      incomes,
	}
	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
func (c *IncomeController) GetIncomesByMonth(ctx *gin.Context) {
	var req model.MonthlyRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	params := model.MonthlyParams{
		Uid:   authPayload.Uid,
		Year:  req.Year,
		Month: req.Month,
	}
	valid := util.IsValidYearMonth(req.Year, req.Month)
	if !valid {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: errors.New("invalid input year or month").Error(),
		})
		return
	}
	total, err := c.IncomeUsecase.GetIncomesByMonth(ctx, params)
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
		Data:    total,
	})
}
