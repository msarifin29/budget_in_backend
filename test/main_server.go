package test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/delivery/controller"
	"github.com/stretchr/testify/assert"
)

func NewTestServer(t *testing.T) *controller.Server {
	log := config.NewLogger()
	con, err := config.LoadConfig("..", "dev")
	assert.NoError(t, err)
	server, sErr := controller.NewServer(log, con, "dev")
	assert.NoError(t, sErr)
	return server
}

func TruncateUser(db *sql.DB) {
	db.Exec("TRUNCATE users")
}
