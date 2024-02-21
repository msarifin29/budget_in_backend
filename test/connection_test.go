package test

import (
	"testing"

	"github.com/msarifin29/be_budget_in/internal/config"
)

func TestConnection(t *testing.T) {
	db := config.Connection()

	db.Ping()
}
