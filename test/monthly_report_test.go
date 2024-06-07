package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetMonthlyReportSuccess(t *testing.T) {
	router := NewTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/monthly_report/3dafa83b-ce13-4bda-883b-191f122a76f8", nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "3dafa83b-ce13-4bda-883b-191f122a76f8", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("monthly_report :", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
func TestGetMonthlyReportDetailSuccess(t *testing.T) {
	router := NewTestServer(t)
	param := model.RequestMonthlyReportDetail{Month: "2024-06"}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/monthly-report-detail/", nil)
	q := req.URL.Query()
	q.Add("month", fmt.Sprintf("%v", param.Month))
	req.URL.RawQuery = q.Encode()
	SetAuthorization(t, req, router.TokenMaker, "bearer", "samsul testing", "3dafa83b-ce13-4bda-883b-191f122a76f8", time.Minute)
	router.Engine.ServeHTTP(w, req)
	bytes, err := io.ReadAll(w.Body)
	fmt.Println("detail :", string(bytes))
	assert.Nil(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success", res["message"])
}
