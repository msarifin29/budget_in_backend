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
		UserName: "testing",
		Email:    "asamsul474@gmail.com",
		Password: "123456",
		TypeUser: "personal",
		Balance:  200000,
		Cash:     50000,
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
		UserName: "testing",
		Email:    "jaya@mail.com",
		Password: "123456",
		TypeUser: "personal",
		Balance:  20000,
		Cash:     20000,
		Currency: "IDR",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(body)))
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
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
		Email:    "jayaing@mail.com",
		Password: "123456",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println(string(bytes))
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
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	fmt.Println("body :", string(bytes))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestLoginUserWithUserNonActive(t *testing.T) {
	router := NewTestServer(t)

	params := model.LoginUserRequest{
		Email:    "testing@mail.com",
		Password: "123456",
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
		Uid:      "deb3823d-5581-4e98-896c-06e5aa3bac4a",
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
		Uid: "b9beed09-e6bb-403d-ad3b-cb6560fa2dba",
	}

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/user/%s", user.Uid)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "jaya", "b9beed09-e6bb-403d-ad3b-cb6560fa2dba", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("account :", string(bytes))
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

func TestDeleteUserSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.NonActiveUserRequest{
		Uid: "ea1b7ff6-e56b-4786-a585-e2ef07520370",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/user/delete", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "ea1b7ff6-e56b-4786-a585-e2ef07520370", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("response : ", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
func TestResetPasswordSuccess(t *testing.T) {
	router := NewTestServer(t)
	params := model.ResetPasswordRequest{
		Uid:         "da063cef-9f52-46da-b98f-0c0067e5869d",
		OldPassword: "112233",
		NewPassword: "password",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/user/reset_password", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "da063cef-9f52-46da-b98f-0c0067e5869d", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
func TestResetPasswordFailed(t *testing.T) {
	router := NewTestServer(t)
	params := model.ResetPasswordRequest{
		Uid:         "da063cef-9f52-46da-b98f-0c0067e5869d",
		OldPassword: "112233",
		NewPassword: "password",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/user/reset_password", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "da063cef-9f52-46da-b98f-0c0067e5869d", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestResetPasswordNoAuthorization(t *testing.T) {
	router := NewTestServer(t)
	params := model.ResetPasswordRequest{
		Uid:         "da063cef-9f52-46da-b98f-0c0067e5869d",
		OldPassword: "112233",
		NewPassword: "password",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/user/reset_password", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestCheckEmailIsExistSuccess(t *testing.T) {
	router := NewTestServer(t)
	params := model.CheckEmail{Email: "test@gmail.com"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/check-email/"+params.Email, nil)

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestCheckEmailIsExistFailed(t *testing.T) {
	router := NewTestServer(t)
	params := model.CheckEmail{Email: "asamsul474@gmail.com"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/check-email/"+params.Email, nil)

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body :", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
