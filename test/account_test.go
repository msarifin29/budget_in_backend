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
		Uid:       "b9beed09-e6bb-403d-ad3b-cb6560fa2dba",
		AccountId: "8dc08329-f468-4532-b942-52301d3cd1c2",
		MaxBudget: 20000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/accounts/update_max_budget", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "jaya", "b9beed09-e6bb-403d-ad3b-cb6560fa2dba", time.Minute)
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
		Uid:       "94c7f7d6-3082-414d-a9af-fd3b8bcbeb92",
		AccountId: "29c08f49-f3b3-46ae-b8e1-1686fff0dbf7",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/accounts/max_budget/", nil)

	q := req.URL.Query()
	q.Add("uid", fmt.Sprintf("%v", params.Uid))
	q.Add("account_id", fmt.Sprintf("%v", params.AccountId))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "94c7f7d6-3082-414d-a9af-fd3b8bcbeb92", time.Minute)
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
