package repository

import (
	"github.com/google/uuid"
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"time"
)

// Common errors
var (
	ErrCategoryHasProducts = NewError("category has products")
	ErrCategoryHasChildren = NewError("category has child categories")
)

func NewError(message string) error {
	return &Error{Message: message}
}
func (e *Error) Error() string {
	return e.Message
}

// Custom error type
type Error struct {
	Message string
}

type categoryRepository struct {
	db *database.DB
}

func (c categoryRepository) Create(category *model.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, description, parent_id, is_active, display_order, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := c.db.Exec(
		query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.IsActive,
		category.DisplayOrder,
		category.CreatedAt,
		category.UpdatedAt,
		false,
	)

	return err
}
func (r *categoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category

	query := `SELECT * FROM categories WHERE id = $1`
	err := r.db.Get(&category, query, id)
	if err != nil {
		return nil, err
	}

	// Get product count
	count, err := r.GetProductCount(id)
	if err != nil {
		return nil, err
	}
	category.ProductCount = count

	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*model.Category, error) {
	var category model.Category

	query := `SELECT * FROM categories WHERE slug = $1`
	err := r.db.Get(&category, query, slug)
	if err != nil {
		return nil, err
	}

	// Get product count
	count, err := r.GetProductCount(category.ID)
	if err != nil {
		return nil, err
	}
	category.ProductCount = count

	return &category, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	category.UpdatedAt = time.Now()

	query := `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, parent_id = $4, is_active = $5, display_order = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.IsActive,
		category.DisplayOrder,
		category.UpdatedAt,
		category.ID,
	)

	return err
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	// Check if category has products
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM product_categories WHERE category_id = $1", id)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrCategoryHasProducts
	}

	// Check if category has children
	err = r.db.Get(&count, "SELECT COUNT(*) FROM categories WHERE parent_id = $1", id)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrCategoryHasChildren
	}

	// Delete category
	_, err = r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}

func (r *categoryRepository) List() ([]model.Category, error) {
	var categories []model.Category

	query := `SELECT * FROM categories ORDER BY display_order ASC, name ASC`
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	// Get product count for each category
	for i := range categories {
		count, err := r.GetProductCount(categories[i].ID)
		if err != nil {
			return nil, err
		}
		categories[i].ProductCount = count
	}

	return categories, nil
}

func (r *categoryRepository) GetProductCount(categoryID uuid.UUID) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM product_categories WHERE category_id = $1`
	err := r.db.Get(&count, query, categoryID)

	return count, err
}
func NewCategoryRepository(db *database.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
