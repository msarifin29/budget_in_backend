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

func TestCreateAccountSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateAccountRequest{
		UserId:      "fadab647-cf23-46fc-bd4d-e7d06d32d753",
		AccountName: "",
		Balance:     20000,
		Cash:        20000,
		Currency:    "IDR",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/accounts/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "fadab647-cf23-46fc-bd4d-e7d06d32d753", time.Minute)
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
func TestCreateAccountFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateAccountRequest{
		UserId:      "fadab647-cf23",
		AccountName: "",
		Balance:     20000,
		Cash:        20000,
		Currency:    "IDR",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/accounts/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "fadab647-cf23-46fc-bd4d-e7d06d32d753", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestUpdateMaxBudgetSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateMaxBudgetRequest{
		Uid:       "3dafa83b-ce13-4bda-883b-191f122a76f8",
		AccountId: "155c136d-cddb-4b07-8a29-4a979387ea41",
		MaxBudget: 200000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/accounts/update_max_budget", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "3dafa83b-ce13-4bda-883b-191f122a76f8", time.Minute)
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
func TestGetMaxBudgetSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetMaxBudgetRequest{
		Uid:       "3dafa83b-ce13-4bda-883b-191f122a76f8",
		AccountId: "155c136d-cddb-4b07-8a29-4a979387ea41",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/accounts/max_budget/", nil)

	q := req.URL.Query()
	q.Add("uid", fmt.Sprintf("%v", params.Uid))
	q.Add("account_id", fmt.Sprintf("%v", params.AccountId))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "3dafa83b-ce13-4bda-883b-191f122a76f8", time.Minute)
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
