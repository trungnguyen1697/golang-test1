package repository

import (
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"time"

	"github.com/google/uuid"
)

type wishlistRepository struct {
	db *database.DB
}

func NewWishlistRepository(db *database.DB) WishlistRepository {
	return &wishlistRepository{
		db: db,
	}
}

type DashboardService struct {
	categoryRepo CategoryRepository
	productRepo  ProductRepository
	reviewRepo   ReviewRepository
	wishlistRepo WishlistRepository
}

func NewDashboardService(
	categoryRepo CategoryRepository,
	productRepo ProductRepository,
	reviewRepo ReviewRepository,
	wishlistRepo WishlistRepository,
) *DashboardService {
	return &DashboardService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
		reviewRepo:   reviewRepo,
		wishlistRepo: wishlistRepo,
	}
}
func (s *DashboardService) GetDashboardStats() (*model.DashboardStats, error) {
	// Get categories with product counts
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	// Create category distribution map
	categoryDistribution := make(map[string]int)
	for _, category := range categories {
		categoryDistribution[category.Name] = category.ProductCount
	}

	// Get recent reviews
	reviewsPage, reviewCount, err := s.reviewRepo.List(0, 5)
	if err != nil {
		return nil, err
	}

	// Build dashboard stats
	stats := &model.DashboardStats{
		CategoryDistribution: categoryDistribution,
		ReviewCount:          reviewCount,
		RecentReviews:        reviewsPage,
		LastUpdated:          time.Now(),
	}

	return stats, nil
}
func (r *wishlistRepository) Add(userID, productID uuid.UUID) error {
	// Check if item already exists
	exists, err := r.IsInWishlist(userID, productID)
	if err != nil {
		return err
	}

	if exists {
		return nil // Already in wishlist
	}

	// Add to wishlist
	query := `INSERT INTO wishlist (user_id, product_id, added_at) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(query, userID, productID, time.Now())

	return err
}

func (r *wishlistRepository) Remove(userID, productID uuid.UUID) error {
	query := `DELETE FROM wishlist WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.Exec(query, userID, productID)
	return err
}

func (r *wishlistRepository) GetByUserID(userID uuid.UUID) ([]model.WishlistItem, error) {
	var items []model.WishlistItem

	query := `
		SELECT w.user_id, w.product_id, w.added_at,
			p.id, p.name, p.description, p.price, p.stock_quantity, p.status
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		WHERE w.user_id = $1
		ORDER BY w.added_at DESC
	`

	rows, err := r.db.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.WishlistItem
		var product model.Product

		err := rows.Scan(
			&item.UserID, &item.ProductID, &item.AddedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.StockQuantity, &product.Status,
		)

		if err != nil {
			return nil, err
		}

		item.Product = &product
		items = append(items, item)
	}

	return items, nil
}

func (r *wishlistRepository) IsInWishlist(userID, productID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM wishlist WHERE user_id = $1 AND product_id = $2)`
	err := r.db.Get(&exists, query, userID, productID)
	return exists, err
}

func (r *wishlistRepository) GetDashboardStats(userID uuid.UUID) (*model.DashboardStats, error) {
	// Create dashboard stats object
	stats := &model.DashboardStats{
		CategoryDistribution: make(map[string]int),
		LastUpdated:          time.Now(),
	}

	// Get wishlist count
	countQuery := `SELECT COUNT(*) FROM wishlist WHERE user_id = $1`
	err := r.db.Get(&stats.WishlistCount, countQuery, userID)
	if err != nil {
		return nil, err
	}

	// Get review count
	reviewCountQuery := `SELECT COUNT(*) FROM reviews WHERE user_id = $1 AND is_deleted = FALSE`
	err = r.db.Get(&stats.ReviewCount, reviewCountQuery, userID)
	if err != nil {
		return nil, err
	}

	// Get category distribution for wishlist items
	categoryQuery := `
		SELECT c.name, COUNT(w.product_id) as count
		FROM categories c
		JOIN product_categories pc ON c.id = pc.category_id
		JOIN wishlist w ON pc.product_id = w.product_id
		WHERE w.user_id = $1
		GROUP BY c.name
	`

	rows, err := r.db.Queryx(categoryQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, err
		}
		stats.CategoryDistribution[category] = count
	}

	// Get top rated products (limit to 5)
	topProductsQuery := `
		SELECT p.id, p.name, p.description, p.price, p.stock_quantity, p.status, 
		       AVG(r.rating) as avg_rating, COUNT(r.id) as review_count
		FROM products p
		JOIN reviews r ON p.id = r.product_id
		WHERE r.is_deleted = FALSE
		GROUP BY p.id
		ORDER BY avg_rating DESC, review_count DESC
		LIMIT 5
	`

	topProductsRows, err := r.db.Queryx(topProductsQuery)
	if err != nil {
		return nil, err
	}
	defer topProductsRows.Close()

	for topProductsRows.Next() {
		var product model.Product
		var avgRating float64
		var reviewCount int

		if err := topProductsRows.Scan(
			&product.ID, &product.Name, &product.Description,
			&product.Price, &product.StockQuantity, &product.Status,
			&avgRating, &reviewCount,
		); err != nil {
			return nil, err
		}

		stats.TopRatedProducts = append(stats.TopRatedProducts, product)
	}

	// Get recent reviews by the user
	recentReviewsQuery := `
		SELECT r.id, r.product_id, r.user_id, r.rating, r.title, r.comment, 
		       r.is_verified_purchase, r.helpful_votes, r.created_at, r.updated_at,
		       u.username, p.name as product_name
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN products p ON r.product_id = p.id
		WHERE r.user_id = $1 AND r.is_deleted = FALSE
		ORDER BY r.created_at DESC
		LIMIT 5
	`

	recentReviewsRows, err := r.db.Queryx(recentReviewsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer recentReviewsRows.Close()

	for recentReviewsRows.Next() {
		var review model.Review
		if err := recentReviewsRows.StructScan(&review); err != nil {
			return nil, err
		}
		stats.RecentReviews = append(stats.RecentReviews, review)
	}

	return stats, nil
}
