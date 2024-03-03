package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type CreditRepository interface {
	CreateCredit(ctx context.Context, tx *sql.Tx, credit model.Credit) (model.Credit, error)
	UpdateCredit(ctx context.Context, tx *sql.Tx, credit model.UpdateCreditRequest) (bool, error)
	CreateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.HistoryCredit) (model.HistoryCredit, error)
	UpdateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.UpdateHistoryCreditRequest) (bool, error)
}

type CreditRepositoryImpl struct {
}

// CreateHistoryCredit implements CreditRepository.
func (CreditRepositoryImpl) CreateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.HistoryCredit) (model.HistoryCredit, error) {
	script := `insert into history_credit (credit_id,th,total,status,type_payment,payment_time) values (?,?,?,?,?,?)`
	result, errC := tx.ExecContext(ctx, script,
		historyC.CreditId,
		historyC.Th,
		historyC.Total,
		historyC.Status,
		historyC.TypePayment,
		historyC.PaymentTime)

	if errC != nil {
		return model.HistoryCredit{}, errC
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return model.HistoryCredit{}, err
	}

	historyC.Id = float64(lastId)
	return historyC, nil
}

// UpdateCredit implements CreditRepository.
func (CreditRepositoryImpl) UpdateCredit(ctx context.Context, tx *sql.Tx, credit model.UpdateCreditRequest) (bool, error) {
	script := `update credits set status_credit = ? where id = ?`
	_, err := tx.ExecContext(ctx, script, credit.StatusCredit, credit.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateHistoryCredit implements CreditRepository.
func (CreditRepositoryImpl) UpdateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.UpdateHistoryCreditRequest) (bool, error) {
	script := `update credits set status = ? where id = ?`
	_, err := tx.ExecContext(ctx, script, historyC.Status, historyC.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create implements CreditRepository.
func (CreditRepositoryImpl) CreateCredit(ctx context.Context, tx *sql.Tx, credit model.Credit) (model.Credit, error) {
	script := `insert into credits (uid,category_credit,type_credit,total,loan_term,status_credit,installment,payment_time) values (?,?,?,?,?,?,?,?)`
	result, errC := tx.ExecContext(ctx, script,
		credit.Uid,
		credit.CategoryCredit,
		credit.TypeCredit,
		credit.Total,
		credit.LoanTerm,
		credit.StatusCredit,
		credit.Installment,
		credit.PaymentTime)
	if errC != nil {
		return model.Credit{}, errC
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return model.Credit{}, err
	}
	credit.Id = float64(lastId)
	return credit, nil
}

func NewCreditRepository() CreditRepository {
	return CreditRepositoryImpl{}
}
