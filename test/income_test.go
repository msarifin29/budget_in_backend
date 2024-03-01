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
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryIncome: util.DAILY,
		TypeIncome:     util.CASH,
		Total:          2000,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/incomes/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
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
		CategoryIncome: util.DAILY,
		Page:           1,
		TotalPage:      10,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/incomes/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("category_income", fmt.Sprintf("%v", params.CategoryIncome))
	q.Add("page", fmt.Sprintf("%d", params.Page))
	q.Add("total_page", fmt.Sprintf("%d", params.TotalPage))
	req.URL.RawQuery = q.Encode()

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
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
