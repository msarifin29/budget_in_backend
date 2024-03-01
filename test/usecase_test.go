package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/msarifin29/be_budget_in/internal/config"
	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/internal/repository"
	"github.com/msarifin29/be_budget_in/internal/usecase"
	"github.com/msarifin29/be_budget_in/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenses(t *testing.T) {
	log := config.NewLogger()
	db := config.Connection(log)
	repoBalance := repository.NewBalanceRepository()
	repoEx := repository.NewExpenseRepository()
	use := usecase.NewExpenseUsecase(repoEx, repoBalance, log, db)
	params := model.CreateExpenseRequest{
		Uid:         "deb3823d-5581-4e98-896c-06e5aa3bac4a",
		ExpenseType: util.DEBIT,
		Category:    util.OTHER,
		Total:       2500,
	}
	value, err := use.CreateExpense(context.Background(), params)
	assert.NoError(t, err)
	fmt.Println("value => ", value)
}
