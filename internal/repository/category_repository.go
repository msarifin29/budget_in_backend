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
	CreateCategoryCredits(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error)
	GetCategoryCredit(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error)
}

type CategoryRepositoryImpl struct{}

// CreateCategoryCredits implements CategoryRepository.
func (CategoryRepositoryImpl) CreateCategoryCredits(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error) {
	script := `insert into t_category_credits (category_id,id,title) values (?,?,?)`
	_, err := tx.ExecContext(ctx, script, category.CategoryId, category.Id, category.Title)
	return category, err
}

// CreateCategoryExpense implements CategoryRepository.
func (CategoryRepositoryImpl) CreateCategoryExpense(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error) {
	script := `insert into t_category_expenses (category_id,id,title) values (?,?,?)`
	_, err := tx.ExecContext(ctx, script, category.CategoryId, category.Id, category.Title)
	return category, err
}

// CreateCategoryIncomes implements CategoryRepository.
func (CategoryRepositoryImpl) CreateCategoryIncomes(ctx context.Context, tx *sql.Tx, category model.Category) (model.Category, error) {
	script := `insert into t_category_incomes (category_id,id,title) values (?,?,?)`
	_, err := tx.ExecContext(ctx, script, category.CategoryId, category.Id, category.Title)
	return category, err
}

// GetCategoryCredit implements CategoryRepository.
func (CategoryRepositoryImpl) GetCategoryCredit(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error) {
	script := `select category_id, id, title from t_category_credits where category_id = ?`
	row := tx.QueryRowContext(ctx, script, categoryId)
	var category model.Category
	err := row.Scan(&category.CategoryId, &category.Id, &category.Title)
	return category, err
}

// GetCategoryExpense implements CategoryRepository.
func (CategoryRepositoryImpl) GetCategoryExpense(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error) {
	script := `select category_id, id, title from t_category_expenses where category_id = ?`
	row := tx.QueryRowContext(ctx, script, categoryId)
	var category model.Category
	err := row.Scan(&category.CategoryId, &category.Id, &category.Title)
	return category, err
}

// GetCategoryIncome implements CategoryRepository.
func (CategoryRepositoryImpl) GetCategoryIncome(ctx context.Context, tx *sql.Tx, categoryId float64) (model.Category, error) {
	script := `select category_id, id, title from t_category_incomes where category_id = ?`
	row := tx.QueryRowContext(ctx, script, categoryId)
	var category model.Category
	err := row.Scan(&category.CategoryId, &category.Id, &category.Title)
	return category, err
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepositoryImpl{}
}
