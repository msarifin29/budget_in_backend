package test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	db := config.Connection(config.NewLogger(), "dev")

	err := db.Ping()
	assert.NoErrorf(t, err, "error %t")
}
