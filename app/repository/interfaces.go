package repository

import (
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	NewUserRepo(db *database.DB) UserRepository
	Exists(username, email string) (bool, error)
	Create(b *model.RegisterUser) error
	All(limit int, offset uint) ([]*model.User, error)
	Get(ID string) (*model.User, error)
	Update(ID string, user *model.UpdateUser) error
	Delete(ID string) error
	ChangePassword(ID string, newPassword string) error
	GetPassword(ID string) (string, error)
}

type ProductRepository interface {
	Create(product *model.Product, categoryIDs []uuid.UUID) error
	GetByID(id uuid.UUID) (*model.Product, error)
	Update(product *model.Product, categoryIDs []uuid.UUID) error
	Delete(id uuid.UUID) error
	List(offset, limit int, search string, categoryID *uuid.UUID, status string) ([]model.Product, int, error)
	ListWithFilters(
		offset, limit int,
		search string,
		categoryID *uuid.UUID,
		status string,
		minPrice, maxPrice *float64,
		minStock, maxStock *int,
		createdAfter, createdBefore *time.Time,
		sortBy, sortOrder string,
	) ([]model.Product, int, error)
	GetCategories(productID uuid.UUID) ([]model.Category, error)
}

type CategoryRepository interface {
	Create(category *model.Category) error
	GetByID(id uuid.UUID) (*model.Category, error)
	GetBySlug(slug string) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	List() ([]model.Category, error)
	GetProductCount(categoryID uuid.UUID) (int, error)
}
type WishlistRepository interface {
	Add(userID, productID uuid.UUID) error
	Remove(userID, productID uuid.UUID) error
	GetByUserID(userID uuid.UUID) ([]model.WishlistItem, error)
	IsInWishlist(userID, productID uuid.UUID) (bool, error)
	GetDashboardStats(userID uuid.UUID) (*model.DashboardStats, error)
}
type ReviewRepository interface {
	Create(review *model.Review) error
	GetByID(id uuid.UUID) (*model.Review, error)
	GetByProductID(productID uuid.UUID) ([]model.Review, error)
	GetByUserID(userID uuid.UUID) ([]model.Review, error)
	Update(review *model.Review) error
	Delete(id uuid.UUID) error
	List(offset, limit int) ([]model.Review, int, error)
}
