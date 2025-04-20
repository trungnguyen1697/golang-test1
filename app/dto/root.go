package dto

import (
	"golang-test1/app/model"
	"time"

	"github.com/google/uuid"
)

// ###########################
// ## Data Transfer Objects ##
// ###########################

type User struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsActive  bool       `json:"is_active"`
	Role      string     `json:"role"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
}

func ToUser(u *model.User) *User {
	return &User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: &u.UpdatedAt,
		IsActive:  u.IsActive,
		Role:      u.Role,
		UserName:  u.UserName,
		Email:     u.Email,
		FullName:  u.FullName,
	}
}

func ToUsers(users []*model.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = ToUser(user)
	}
	return res
}

// Category DTO
type Category struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Slug         string     `json:"slug"`
	Description  string     `json:"description,omitempty"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	DisplayOrder int        `json:"display_order"`
	ProductCount int        `json:"product_count,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func ToCategory(c *model.Category) *Category {
	return &Category{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		ParentID:     c.ParentID,
		IsActive:     c.IsActive,
		DisplayOrder: c.DisplayOrder,
		ProductCount: c.ProductCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func ToCategories(categories []model.Category) []*Category {
	res := make([]*Category, len(categories))
	for i, category := range categories {
		res[i] = ToCategory(&category)
	}
	return res
}

// Product DTO
type Product struct {
	ID            uuid.UUID      `json:"id"`
	SKU           string         `json:"sku"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Price         float64        `json:"price"`
	SalePrice     float64        `json:"sale_price,omitempty"`
	CostPrice     float64        `json:"cost_price,omitempty"`
	StockQuantity int            `json:"stock_quantity"`
	Status        string         `json:"status"`
	Attributes    map[string]any `json:"attributes,omitempty"`
	Categories    []*Category    `json:"categories,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func ToProduct(p *model.Product) *Product {
	return &Product{
		ID:            p.ID,
		SKU:           p.SKU,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		SalePrice:     p.SalePrice,
		CostPrice:     p.CostPrice,
		StockQuantity: p.StockQuantity,
		Status:        p.Status,
		Attributes:    p.Attributes,
		Categories:    ToCategories(p.Categories),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func ToProducts(products []model.Product) []*Product {
	res := make([]*Product, len(products))
	for i, product := range products {
		product := product // Create a new variable to avoid issues with the closure
		res[i] = ToProduct(&product)
	}
	return res
}

// WishlistItem DTO for adding a product to wishlist
type WishlistItemRequest struct {
	ProductID string `json:"product_id" example:"5c9f8f9e-7c1f-4b9c-8c1f-7c1f4b9c8c1f"`
}
