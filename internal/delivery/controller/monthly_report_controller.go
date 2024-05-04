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

type MonthlyReportController struct {
	Usecase usecase.MonthlyReportUsecase
	Log     *logrus.Logger
}

func NewMonthlyController(Usecase usecase.MonthlyReportUsecase, Log *logrus.Logger) *MonthlyReportController {
	return &MonthlyReportController{Usecase: Usecase, Log: Log}
}
func (c *MonthlyReportController) GetMonthlyReport(ctx *gin.Context) {
	var req model.ParamMonthlyReport
	if err := ctx.ShouldBindUri(&req); err != nil {
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
	res, err := c.Usecase.GetMonthlyReport(ctx, req)
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
func (c *MonthlyReportController) GetMonthlyReportDetail(ctx *gin.Context) {
	var req model.RequestMonthlyReportDetail
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Errorf("failed binding request query with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	records, err := c.Usecase.GetMonthlyReportDetail(ctx, model.ParamMonthlyReportDetail{
		Uid:   authPayload.Uid,
		Month: req.Month,
	})
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
		Data:    records,
	})
}
