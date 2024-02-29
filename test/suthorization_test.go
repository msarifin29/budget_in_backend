package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	delivery "github.com/msarifin29/be_budget_in/internal/delivery/middleware"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/stretchr/testify/assert"
)

func SetAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker util.Maker,
	authorizationType string,
	username string,
	uid string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, duration, uid)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("authorization", authorizationHeader)
}

func TestAuthorizationSuccess(t *testing.T) {
	server := NewTestServer(t)
	server.Engine.GET("/api/auth", delivery.AuthMiddleware(server.TokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/auth", nil)
	assert.NoError(t, err)
	SetAuthorization(t, request, server.TokenMaker, "bearer", "samsul", "fadab647-cf23-46fc-bd4d-e7d06d32d753", time.Minute)
	server.Engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
func TestNoAuthorization(t *testing.T) {
	server := NewTestServer(t)
	server.Engine.GET("/api/auth", delivery.AuthMiddleware(server.TokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/auth", nil)
	assert.NoError(t, err)
	server.Engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
func TestInvalidAuthorization(t *testing.T) {
	server := NewTestServer(t)
	server.Engine.GET("/api/auth", delivery.AuthMiddleware(server.TokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/auth", nil)
	assert.NoError(t, err)
	SetAuthorization(t, request, server.TokenMaker, "", "samsul", "fadab647-cf23-46fc-bd4d-e7d06d32d753", time.Minute)
	server.Engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
func TestExpiredAuthorization(t *testing.T) {
	server := NewTestServer(t)
	server.Engine.GET("/api/auth", delivery.AuthMiddleware(server.TokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/api/auth", nil)
	assert.NoError(t, err)
	SetAuthorization(t, request, server.TokenMaker, "bearer", "samsul", "fadab647-cf23-46fc-bd4d-e7d06d32d753", -time.Minute)
	server.Engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
