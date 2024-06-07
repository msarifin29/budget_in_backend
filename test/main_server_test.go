package test

import (
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/delivery/controller"
	"github.com/stretchr/testify/assert"
)

func TestMainServer(t *testing.T) {
	log := config.NewLogger()
	con, err := config.LoadConfig("..", "prod")
	assert.NoError(t, err)
	server, sErr := controller.NewServer(log, con)
	assert.NoError(t, sErr)
	err = server.Start(con.ServerAddress)
	assert.NoError(t, err)
	time.Sleep(time.Duration(5))
	os.Exit(1)
}
