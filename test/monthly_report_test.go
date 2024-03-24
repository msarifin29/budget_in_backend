package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMonthlyReportSuccess(t *testing.T) {
	router := NewTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/user/monthly_report/b9beed09-e6bb-403d-ad3b-cb6560fa2dba", nil)

	SetAuthorization(t, req, router.TokenMaker, "bearer", "jaya", "b9beed09-e6bb-403d-ad3b-cb6560fa2dba", time.Minute)
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
