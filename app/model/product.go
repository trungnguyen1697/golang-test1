package model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	SKU           string         `json:"sku" db:"sku"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Price         float64        `json:"price" db:"price"`
	SalePrice     float64        `json:"sale_price,omitempty" db:"sale_price"`
	CostPrice     float64        `json:"cost_price,omitempty" db:"cost_price"`
	StockQuantity int            `json:"stock_quantity" db:"stock_quantity"`
	Status        string         `json:"status" db:"status"`
	Attributes    map[string]any `json:"attributes,omitempty" db:"attributes"`
	Categories    []Category     `json:"categories,omitempty"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	IsDeleted     bool           `json:"is_deleted,omitempty" db:"is_deleted"`
}

// Product creation/update input
type ProductInput struct {
	Name          string         `json:"name" validate:"required"`
	SKU           string         `json:"sku" validate:"required"`
	Description   string         `json:"description"`
	Price         float64        `json:"price" validate:"required,gt=0"`
	SalePrice     *float64       `json:"sale_price,omitempty"`
	CostPrice     *float64       `json:"cost_price,omitempty"`
	StockQuantity int            `json:"stock_quantity" validate:"required,gte=0"`
	Status        string         `json:"status" validate:"required,oneof=active inactive out_of_stock"`
	Attributes    map[string]any `json:"attributes,omitempty"`
	CategoryIDs   []uuid.UUID    `json:"category_ids" validate:"required,min=1"`
}
