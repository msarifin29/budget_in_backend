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
		Installment:    3000,
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
