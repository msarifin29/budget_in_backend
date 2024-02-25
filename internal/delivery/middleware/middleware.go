package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util"
)

const AuthorizationHeaderKey = "authorization"
const AuthorizationTypeBearer = "bearer"
const AuthorizationPayloadKey = "authorization_payload"

func AuthMiddleware(tokenMaker util.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provide")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.MetaErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.MetaErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.MetaErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.MetaErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
