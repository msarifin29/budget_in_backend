package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateUserRequest{
		UserName: "samsul",
		Email:    "samsul@mail.com",
		Password: "123456",
		TypeUser: "personal",
		Balance:  20000,
		Savings:  20000,
		Cash:     20000,
		Debts:    0,
		Currency: "IDR",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
	assert.NotEmpty(t, res["data"])
}

func TestCreateUserDuplicate(t *testing.T) {
	router := NewTestServer(t)
	params := model.CreateUserRequest{
		UserName: "kai",
		Email:    "test@mail.com",
		Password: "123456",
		TypeUser: "personal",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(body)))
	router.Engine.ServeHTTP(w, req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateUserInvalidType(t *testing.T) {
	router := NewTestServer(t)
	params := model.CreateUserRequest{
		UserName: "katakuri",
		Email:    "katakuri@mail.com",
		Password: "123456",
		TypeUser: "type",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(body)))
	router.Engine.ServeHTTP(w, req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateUserRequestError(t *testing.T) {
	router := NewTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", nil)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUserSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "123456",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
	assert.NotEmpty(t, res["data"])
}

func TestLoginUserInvalidEmail(t *testing.T) {
	router := NewTestServer(t)

	params := model.LoginUserRequest{
		Email:    "aa@mail.com",
		Password: "123456",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestLoginUserInvalidPassword(t *testing.T) {
	router := NewTestServer(t)

	params := model.LoginUserRequest{
		Email:    "test@mail.com",
		Password: "password",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateteUserSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateUserRequest{
		UserName: "samsul testing",
		Uid:      "fe317556e-74b2-4199-8a30-33bd56fc5e9e",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
func TestUpdateteUserWithIdInvalid(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateUserRequest{
		UserName: "samsul testing",
		Uid:      "f1687230",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestUnAuthorizationUser(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateUserRequest{
		UserName: "samsul testing",
		Uid:      "f1687230-49d3-4657-96be-9b934ed0387f",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/update", strings.NewReader(string(body)))
	router.Engine.ServeHTTP(w, req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserByIdSuccess(t *testing.T) {
	router := NewTestServer(t)
	user := model.UserRequest{
		Uid: "f1687230-49d3-4657-96be-9b934ed0387f",
	}

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/user/%s", user.Uid)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
func TestGetUserInvalidId(t *testing.T) {
	router := NewTestServer(t)
	user := model.UserRequest{
		Uid: "f1687230",
	}

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/user/%s", user.Uid)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestGetUserNoAuthorization(t *testing.T) {
	router := NewTestServer(t)
	user := model.UserRequest{
		Uid: "f1687230-49d3-4657-96be-9b934ed0387f",
	}

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/user/%s", user.Uid)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
