package controller

import (
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
	var req model.CreateExpenseRequest

	model.BindingResponseError(&req, c.Log, ctx)

	res, err := c.Usecase.CreateExpense(ctx, req)

	model.BadRequestResponseError(err, ctx)

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

	model.NotFoundResponseError(err, ctx)

	ctx.JSON(http.StatusOK, model.MetaResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})
}
func (c *ExpenseController) UpdateExpense(ctx *gin.Context) {
	var req model.UpdateExpenseRequest
	model.BindingResponseError(&req, c.Log, ctx)

	res, err := c.Usecase.UpdateExpense(ctx, req)

	model.BadRequestResponseError(err, ctx)

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

	model.BadRequestResponseError(err, ctx)

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
		Uid:    authPayload.Uid,
		Limit:  req.TotalPage,
		Offset: (req.Page - 1) * req.TotalPage,
	}

	expenses, total, err := c.Usecase.GetExpenses(ctx, params)

	model.BadRequestResponseError(err, ctx)
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
