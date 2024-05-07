package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/templates"
	"github.com/sirupsen/logrus"
)

type PrivacyController struct {
	Log *logrus.Logger
}

func NewPrivacyController(Log *logrus.Logger) *PrivacyController {
	return &PrivacyController{Log: Log}
}

func (c *PrivacyController) PrivacyPolice(ctx *gin.Context) {
	var req model.PrivacyParam

	if err := ctx.ShouldBindUri(&req); err != nil {
		c.Log.Errorf("failed binding request with %t:", err)
		ctx.JSON(http.StatusBadRequest, model.MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if req.Lang == "in" {
		ctx.String(http.StatusOK, templates.PrivacyId)
	} else if req.Lang == "en" {
		ctx.String(http.StatusOK, templates.PrivacyEn)
	} else if req.Lang == "" || req.Lang != "in" || req.Lang == "en" {
		ctx.String(http.StatusBadRequest, "invalid input language")
	} else {
		ctx.String(http.StatusNotFound, "page not found")
	}
}
