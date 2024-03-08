package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/usecase"
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
