package test

import (
	"testing"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadDev(t *testing.T) {
	c, err := config.LoadConfig("..", "env_dev")
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
}
