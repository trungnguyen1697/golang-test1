package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	ProductID          uuid.UUID `json:"product_id" db:"product_id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Rating             int       `json:"rating" db:"rating"`
	Title              string    `json:"title,omitempty" db:"title"`
	Comment            string    `json:"comment" db:"comment"`
	IsVerifiedPurchase bool      `json:"is_verified_purchase" db:"is_verified_purchase"`
	HelpfulVotes       int       `json:"helpful_votes" db:"helpful_votes"`
	Username           string    `json:"username,omitempty"`
	ProductName        string    `json:"product_name,omitempty" db:"product_name"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted          bool      `json:"is_deleted,omitempty" db:"is_deleted"`
}

// Review creation input
type ReviewInput struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Title     string    `json:"title"`
	Comment   string    `json:"comment" validate:"required"`
}
