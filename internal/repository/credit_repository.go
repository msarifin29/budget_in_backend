package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
	"github.com/msarifin29/be_budget_in/util/zero"
)

type CreditRepository interface {
	CreateCredit(ctx context.Context, tx *sql.Tx, credit model.Credit) (model.Credit, error)
	UpdateCredit(ctx context.Context, tx *sql.Tx, credit model.UpdateCreditRequest) (bool, error)
	UpdateTotalCredit(ctx context.Context, tx *sql.Tx, uid string, id float64, total float64) (bool, error)
	GetCreditById(ctx context.Context, tx *sql.Tx, credit model.GetCreditRequest) (model.Credit, error)
	CreateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.HistoryCredit) (model.HistoryCredit, error)
	UpdateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.UpdateHistoryCreditParams) (bool, error)
	GetHistoryCreditById(ctx context.Context, tx *sql.Tx, credit model.GetHistoryCreditRequest) (model.HistoryCredit, error)
}

type CreditRepositoryImpl struct {
}

// UpdateTotalCredit implements CreditRepository.
func (CreditRepositoryImpl) UpdateTotalCredit(ctx context.Context, tx *sql.Tx, uid string, id float64, total float64) (bool, error) {
	script := `update credits set total = ? where uid = ? && id = ?`
	_, err := tx.ExecContext(ctx, script, total, uid, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetCreditById implements CreditRepository.
func (CreditRepositoryImpl) GetCreditById(ctx context.Context, tx *sql.Tx, credit model.GetCreditRequest) (model.Credit, error) {
	script := `select * from credits where uid = ? && id = ?`
	rows := tx.QueryRowContext(ctx, script, credit.Uid, credit.Id)
	var i model.Credit
	update := zero.TimeFromPtr(&i.UpdatedAt)
	err := rows.Scan(
		&i.Uid,
		&i.Id,
		&i.CategoryCredit,
		&i.TypeCredit,
		&i.Total,
		&i.LoanTerm,
		&i.StatusCredit,
		&i.CreatedAt,
		&update,
		&i.Installment,
		&i.PaymentTime,
	)
	return i, err
}

// GetHistoryCreditById implements CreditRepository.
func (CreditRepositoryImpl) GetHistoryCreditById(ctx context.Context, tx *sql.Tx, credit model.GetHistoryCreditRequest) (model.HistoryCredit, error) {
	script := `select * from history_credit where id = ?`
	rows := tx.QueryRowContext(ctx, script, credit.Id)
	var i model.HistoryCredit
	update := zero.TimeFromPtr(&i.UpdatedAt)
	typePayment := zero.StringFromPtr(&i.TypePayment)
	err := rows.Scan(
		&i.CreditId,
		&i.Id,
		&i.Th,
		&i.Total,
		&i.Status,
		&i.CreatedAt,
		&update,
		&typePayment,
		&i.PaymentTime,
	)
	return i, err
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
	script := `update credits set status_credit = ? where uid = ? && id = ?`
	_, err := tx.ExecContext(ctx, script, credit.StatusCredit, credit.Uid, credit.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateHistoryCredit implements CreditRepository.
func (CreditRepositoryImpl) UpdateHistoryCredit(ctx context.Context, tx *sql.Tx, historyC model.UpdateHistoryCreditParams) (bool, error) {
	script := `update history_credit set status = ? , type_payment = ? where id = ? and credit_id = ?`
	_, err := tx.ExecContext(ctx, script, historyC.Status, historyC.TypePayment, historyC.Id, historyC.CreditId)
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
