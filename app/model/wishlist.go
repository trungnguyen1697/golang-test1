package model

import (
	"time"

	"github.com/google/uuid"
)

type WishlistItem struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	AddedAt   time.Time `json:"added_at" db:"added_at"`
	Product   *Product  `json:"product,omitempty"`
}

// Dashboard stats model
type DashboardStats struct {
	CategoryDistribution map[string]int `json:"category_distribution"`
	WishlistCount        int            `json:"wishlist_count"`
	ReviewCount          int            `json:"review_count"`
	TopRatedProducts     []Product      `json:"top_rated_products,omitempty"`
	RecentReviews        []Review       `json:"recent_reviews,omitempty"`
	LastUpdated          time.Time      `json:"last_updated"`
}
