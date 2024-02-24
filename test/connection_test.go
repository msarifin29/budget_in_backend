package test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/msarifin29/be_budget_in/internal/config"
)

func TestConnection(t *testing.T) {
	db := config.Connection(config.NewLogger())

	db.Ping()
}
