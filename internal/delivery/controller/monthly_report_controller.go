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

type MonthlyReportController struct {
	Usecase usecase.MonthlyReportUsecase
	Log     *logrus.Logger
}

func NewMonthlyController(Usecase usecase.MonthlyReportUsecase, Log *logrus.Logger) *MonthlyReportController {
	return &MonthlyReportController{Usecase: Usecase, Log: Log}
}
func (c *MonthlyReportController) GetMonthlyReport(ctx *gin.Context) {

	authPayload := ctx.MustGet(delivery.AuthorizationPayloadKey).(*util.Payload)
	res, err := c.Usecase.GetMonthlyReport(ctx, authPayload.Uid)
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
