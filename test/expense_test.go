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

func TestCreateExpenseSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseRequest{
		Uid:         "f1687230-49d3-4657-96be-9b934ed0387f",
		ExpenseType: "Cash",
		Total:       75000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/expenses/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
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
func TestCreateExpenseFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseRequest{
		ExpenseType: "Cash",
		Total:       45000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/expenses/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateExpenseUnAuthorized(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseRequest{
		Uid:         "f1687230-49d3-4657-96be-9b934ed0387f",
		ExpenseType: "Cash",
		Total:       45000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/expenses/create", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestGetExpenseByIdSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.ExpenseParamWithId{Id: 10}
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/expenses/%v", params.Id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
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
func TestGetExpenseByIdWithInvalidId(t *testing.T) {
	router := NewTestServer(t)

	params := model.ExpenseParamWithId{Id: 1000}
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/expenses/%v", params.Id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestUpdateExpenseSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateExpenseRequest{
		Id:          9,
		ExpenseType: "Cash",
		Total:       90000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
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
func TestUpdateExpenseFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateExpenseRequest{
		Id:          9,
		ExpenseType: "Invalid",
		Total:       75000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestUpdateExpenseUnAuthorized(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateExpenseRequest{
		Id:          9,
		ExpenseType: "Cash",
		Total:       75000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestDeleteExpenseSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.ExpenseParamWithId{Id: 6}
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/expenses/%v", params.Id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestDeleteExpenseFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.ExpenseParamWithId{Id: 1000000}
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/expenses/%v", params.Id)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetExpensesSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetExpenseRequest{
		Uid:       "f1687230-49d3-4657-96be-9b934ed0387f",
		Page:      1,
		TotalPage: 5,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/expenses/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("uid", fmt.Sprintf("%v", params.Uid))
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res model.MetaResponse
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res.Message)
	assert.NotEmpty(t, res.Data)
}
func TestGetExpensesFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetExpenseRequest{
		Uid:       "f1687230-49d3-4657-96be-9b934ed0387f",
		Page:      1,
		TotalPage: 5,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/expenses/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("uid", fmt.Sprintf("%v", params.Uid))
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	var res model.MetaResponse
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res.Message)
	assert.NotEmpty(t, res.Data)
}
