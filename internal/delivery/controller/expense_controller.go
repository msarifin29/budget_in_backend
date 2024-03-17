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

type ExpenseController struct {
	Usecase usecase.ExpenseUsecase
	Log     *logrus.Logger
}

func NewExpenseController(Usecase usecase.ExpenseUsecase, Log *logrus.Logger) *ExpenseController {
	return &ExpenseController{Usecase: Usecase, Log: Log}
}

func (c *ExpenseController) CreateExpense(ctx *gin.Context) {
	var req model.CreateExpenseParams
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
	res, err := c.Usecase.CreateExpense(ctx, req)
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
func (c *ExpenseController) GetExpenseById(ctx *gin.Context) {
	var req model.ExpenseParamWithId
	if err := ctx.ShouldBindUri(&req); err != nil {
		c.Log.Errorf("failed binding request with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	res, err := c.Usecase.GetExpenseById(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusNotFound, model.MetaErrorResponse{
			Code:    http.StatusNotFound,
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
func (c *ExpenseController) UpdateExpense(ctx *gin.Context) {
	var req model.UpdateExpenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Errorf("failed binding request %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	res, err := c.Usecase.UpdateExpense(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK,
		model.MetaResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    res,
		})
}
func (c *ExpenseController) DeleteExpense(ctx *gin.Context) {
	var req model.ExpenseParamWithId
	if err := ctx.ShouldBindUri(&req); err != nil {
		c.Log.Errorf("failed binding request with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	x, err := c.Usecase.GetExpenseById(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid input id",
		})
		return
	}
	err = c.Usecase.DeleteExpense(ctx, x.Id)
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
		Data:    true,
	})
}
func (c *ExpenseController) GetExpenses(ctx *gin.Context) {
	var req model.GetExpenseRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	params := model.GetExpenseParams{
		Uid:         authPayload.Uid,
		Status:      req.Status,
		Category:    req.Category,
		ExpenseType: req.ExpenseType,
		Limit:       req.TotalPage,
		Offset:      (req.Page - 1) * req.TotalPage,
	}

	expenses, total, err := c.Usecase.GetExpenses(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	lastPage := int32(total/float64(req.TotalPage) + 1)
	res := model.ExpensesResponse{
		Page:      req.Page,
		TotalPage: req.TotalPage,
		LastPage:  lastPage,
		Total:     int32(total),
		Data:      expenses,
	}
	ctx.JSON(http.StatusOK,
		model.MetaResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    res,
		})
}
func (c *ExpenseController) GetExpensesByMonth(ctx *gin.Context) {
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

	total, err := c.Usecase.GetExpensesByMonth(ctx, params)
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
