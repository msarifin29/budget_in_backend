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

func TestCreateCreditSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateCreditRequest{
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryCredit: util.ELECTRONIC,
		TypeCredit:     util.MONTHLY,
		LoanTerm:       3,
		Installment:    2500,
		PaymentTime:    time.Now().Day(),
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/credits/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body =>", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
	assert.NotEmpty(t, res["data"])
}
func TestCreateCreditFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateCreditRequest{
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryCredit: util.ELECTRONIC,
		TypeCredit:     util.MONTHLY,
		LoanTerm:       0,
		Installment:    1000,
		PaymentTime:    time.Now().Day(),
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/credits/create", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	fmt.Println("body =>", string(bytes))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestCreateCreditNoAuthorization(t *testing.T) {
	router := NewTestServer(t)

	params := model.CreateCreditRequest{
		Uid:            "f1687230-49d3-4657-96be-9b934ed0387f",
		CategoryCredit: util.ELECTRONIC,
		TypeCredit:     util.MONTHLY,
		LoanTerm:       3,
		Installment:    1000,
		PaymentTime:    time.Now().Day(),
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/credits/create", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	fmt.Println("body =>", string(bytes))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestUpdateHistoryCreditSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateHistoryCreditRequest{
		Uid:         "f1687230-49d3-4657-96be-9b934ed0387f",
		Id:          7,
		CreditId:    3,
		TypePayment: util.CASH,
		AccountId:   "96be-9b934ed0387f",
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/credits/update_history", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("body =>", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
	assert.NotEmpty(t, res["data"])
}
func TestUpdateHistoryCreditFailed(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateHistoryCreditRequest{
		Uid:         "f1687230-49d3-4657-96be-9b934ed0387f",
		Id:          3,
		TypePayment: util.ELECTRONIC,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/credits/update_history", strings.NewReader(string(body)))

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "f1687230-49d3-4657-96be-9b934ed0387f", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	fmt.Println("body =>", string(bytes))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestUpdateHistoryCreditNoAuthorization(t *testing.T) {
	router := NewTestServer(t)

	params := model.UpdateHistoryCreditRequest{
		Uid:         "f1687230-49d3-4657-96be-9b934ed0387f",
		Id:          3,
		TypePayment: util.ELECTRONIC,
	}
	body, err := json.Marshal(params)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/credits/update_history", strings.NewReader(string(body)))

	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	fmt.Println("body =>", string(bytes))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
func TestGetCreditsSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetCreditsRequest{
		Page:      1,
		TotalPage: 5,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/credits/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
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
func TestGetHistoriesCreditsSuccess(t *testing.T) {
	router := NewTestServer(t)

	params := model.GetHistoriesCreditsRequest{
		CreditId:  3,
		Page:      1,
		TotalPage: 5,
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/histories_credits/", nil)
	// Add query parameters to request URL
	q := req.URL.Query()
	q.Add("credit_id", fmt.Sprintf("%v", params.CreditId))
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
