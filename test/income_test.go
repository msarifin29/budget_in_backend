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

func TestCreateIncomeSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateIncomeRequest{
		Uid: "da063cef-9f52-46da-b98f-0c0067e5869d",
		// CategoryIncome: util.SALARY,
		TypeIncome: util.DEBIT,
		Total:      20000,
		CategoryId: 2,
		AccountId:  "faae4ed7-f719-45a5-b259-3e6bf7407ba0",
		// CreatedAt:  "", // 2015-09-02T08:00:00Z
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/incomes/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "da063cef-9f52-46da-b98f-0c0067e5869d", time.Minute)
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
func TestCreateIncomeInvalid(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateIncomeRequest{
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryIncome: util.CANCELLED,
		TypeIncome:     util.CASH,
		Total:          2000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/incomes/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateIncomeNoAuthorization(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateIncomeRequest{
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryIncome: util.WEEKLY,
		TypeIncome:     util.CASH,
		Total:          2000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/incomes/create", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetIncomesSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetIncomeRequest{
		// CategoryIncome: util.DAILY,
		// TypeIncome: "Cash",
		// CategoryId: 1,
		Page:      1,
		TotalPage: 10,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/incomes/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("category_income", fmt.Sprintf("%v", params.CategoryIncome))
	q.Add("category_id", fmt.Sprintf("%v", params.CategoryId))
	q.Add("type_income", fmt.Sprintf("%v", params.TypeIncome))
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul", "da063cef-9f52-46da-b98f-0c0067e5869d", time.Minute)
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
func TestGetIncomesByMonthSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.MonthlyRequest{Year: "2024", Month: "03"}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/incomes/monthly_report/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("year", fmt.Sprintf("%v", params.Year))
	q.Add("month", fmt.Sprintf("%v", params.Month))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "akainu", "deb3823d-5581-4e98-896c-06e5aa3bac4a", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body : ", string(bytes))
	assert.Nil(t, err)
	var res model.MetaResponse
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res.Message)
}

func TestGetIncomesByMonthFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.MonthlyRequest{Year: "2024", Month: "00"}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/incomes/monthly_report/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("year", fmt.Sprintf("%v", params.Year))
	q.Add("month", fmt.Sprintf("%v", params.Month))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "akainu", "deb3823d-5581-4e98-896c-06e5aa3bac4a", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body : ", string(bytes))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
