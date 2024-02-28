package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func BindingResponseError(T any, Log *logrus.Logger, ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(&T); err != nil {
		Log.Errorf("failed binding request %t:", err)
		ctx.JSON(http.StatusBadRequest, MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
}

func BadRequestResponseError(err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, MetaErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
}

func NotFoundResponseError(err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusNotFound, MetaErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
}

func ServerResponseError(err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MetaErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
}
