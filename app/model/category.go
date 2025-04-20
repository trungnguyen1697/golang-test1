package model

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Slug         string     `json:"slug" db:"slug"`
	Description  string     `json:"description,omitempty" db:"description"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	DisplayOrder int        `json:"display_order" db:"display_order"`
	ProductCount int        `json:"product_count,omitempty"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted    bool       `json:"is_deleted,omitempty" db:"is_deleted"`
}

// Category creation/update input
type CategoryInput struct {
	Name         string     `json:"name" validate:"required"`
	Description  string     `json:"description"`
	ParentID     *uuid.UUID `json:"parent_id,omitempty"`
	IsActive     bool       `json:"is_active"`
	DisplayOrder int        `json:"display_order"`
}
