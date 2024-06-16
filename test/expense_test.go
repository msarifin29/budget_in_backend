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
	"github.com/msarifin29/be_budget_in/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseParams{
		Uid:         "8601f262-5c0f-4024-86db-8f4737360180",
		ExpenseType: util.CASH,
		Total:       5000,
		CategoryId:  4,
		Notes:       "",
		AccountId:   "e38418b8-3342-4d10-b7c2-e09e9fc90193",
		// CreatedAt:   "2024-01-02T08:00:00Z",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/expenses/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "8601f262-5c0f-4024-86db-8f4737360180", time.Minute)
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
func TestCreateExpenseFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseRequest{
		ExpenseType: util.CASH,
		Total:       45000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/expenses/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateExpenseUnAuthorized(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateExpenseRequest{
		Uid:         "3dafa83b-ce13-4bda-883b-191f122a76f8",
		ExpenseType: util.CASH,
		Total:       45000,
		Category:    util.OTHER,
		AccountId:   "155c136d-cddb-4b07-8a29-4a979387ea41",
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

	params := model.ExpenseParamWithId{Id: 87}
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/api/expenses/%v", params.Id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.Nil(t, err)
	fmt.Println(string(bytes))
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

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
func TestUpdateExpenseSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateExpenseRequest{
		Id:          1,
		ExpenseType: util.CASH,
		AccountId:   "e38418b8-3342-4d10-b7c2-e09e9fc90193",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "8601f262-5c0f-4024-86db-8f4737360180", time.Minute)
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
		Id:          86,
		ExpenseType: util.DEBIT,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "boy", "a0ea433a-13b2-4414-aa8a-6369acd2b547", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestUpdateExpenseUnAuthorized(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateExpenseRequest{
		Id: 10,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/expenses/update", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// func TestDeleteExpenseSuccess(t *testing.T) {
// 	router := NewTestServer(t)

// 	params := model.ExpenseParamWithId{Id: 7}
// 	w := httptest.NewRecorder()
// 	url := fmt.Sprintf("/api/expenses/%v", params.Id)
// 	req, _ := http.NewRequest(http.MethodDelete, url, nil)

// 	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
// 	router.Engine.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
// func TestDeleteExpenseFailed(t *testing.T) {
// 	router := NewTestServer(t)

// 	params := model.ExpenseParamWithId{Id: 1000000}
// 	w := httptest.NewRecorder()
// 	url := fmt.Sprintf("/api/expenses/%v", params.Id)
// 	req, _ := http.NewRequest(http.MethodDelete, url, nil)

// 	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
// 	router.Engine.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// }

func TestGetExpensesSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetExpenseRequest{
		Status: "success",
		// ExpenseType: util.DEBIT,
		// Id: 4,
		CreatedAt: "2024-06-16",
		Page:      1,
		TotalPage: 10,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/expenses/", nil)
	q := req.URL.Query()
	q.Add("expense_type", fmt.Sprintf("%v", params.ExpenseType))
	q.Add("created_at", fmt.Sprintf("%v", params.CreatedAt))
	q.Add("id", fmt.Sprintf("%v", params.Id))
	q.Add("status", fmt.Sprintf("%v", params.Status))
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "testing", "8601f262-5c0f-4024-86db-8f4737360180", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body : ", string(bytes))
	assert.Nil(t, err)
	var res model.MetaResponse
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res.Message)
	assert.NotEmpty(t, res.Data)
}
func TestGetExpensesInvalidInputTotalPage(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetExpenseRequest{
		Page:      1000,
		TotalPage: 5,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/expenses/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
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
