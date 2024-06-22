package repository

import (
	"context"
	"database/sql"

	"github.com/msarifin29/be_budget_in/internal/model"
)

type CategoryRepository interface {
	CreateCategoryExpense(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error)
	GetCategoryExpense(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error)
	CreateCategoryIncomes(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error)
	GetCategoryIncome(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error)
	UpdateStatusCategoryExpense(ctx context.Context, tx *sql.Tx, categoryId float64) error
}

type CategoryRepositoryImpl struct{}

// UpdateStatusCategoryExpense implements CategoryRepository.
func (c CategoryRepositoryImpl) UpdateStatusCategoryExpense(ctx context.Context, tx *sql.Tx, categoryId float64) error {
	script := `update t_category_expenses set status = 'cancelled' where category_id = $1`
	_, err := tx.ExecContext(ctx, script, categoryId)
	return err
}

// CreateCategoryExpense implements CategoryRepository.
func (CategoryRepositoryImpl) CreateCategoryExpense(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error) {
	script := `insert into t_category_expenses (category_id,id,title,user_id) values ($1,$2,$3,$4)`
	_, err := tx.ExecContext(ctx, script, category.CategoryId, category.Id, category.Title, category.UserId)
	return category, err
}

// CreateCategoryIncomes implements CategoryRepository.
func (CategoryRepositoryImpl) CreateCategoryIncomes(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error) {
	script := `insert into t_category_incomes (category_id,id,title,user_id) values ($1,$2,$3,$4)`
	_, err := tx.ExecContext(ctx, script, category.CategoryId, category.Id, category.Title, category.UserId)
	return category, err
}

// GetCategoryExpense implements CategoryRepository.
func (CategoryRepositoryImpl) GetCategoryExpense(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error) {
	script := `select category_id, id, title, user_id from t_category_expenses where category_id = $1`
	row := tx.QueryRowContext(ctx, script, categoryId)
	var category model.Category
	err := row.Scan(&category.CategoryId, &category.Id, &category.Title, &category.UserId)
	return category, err
}

// GetCategoryIncome implements CategoryRepository.
func (CategoryRepositoryImpl) GetCategoryIncome(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error) {
	script := `select category_id, id, title, user_id from t_category_incomes where category_id = $1`
	row := tx.QueryRowContext(ctx, script, categoryId)
	var category model.Category
	err := row.Scan(&category.CategoryId, &category.Id, &category.Title, &category.UserId)
	return category, err
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepositoryImpl{}
}
