package repository

import (
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"time"

	"github.com/google/uuid"
)

type reviewRepository struct {
	db *database.DB
}

func NewReviewRepository(db *database.DB) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}
func (r *reviewRepository) Create(review *model.Review) error {
	query := `
		INSERT INTO reviews (id, product_id, user_id, rating, title, comment, 
                           is_verified_purchase, helpful_votes, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(
		query,
		review.ID,
		review.ProductID,
		review.UserID,
		review.Rating,
		review.Title,
		review.Comment,
		review.IsVerifiedPurchase,
		review.HelpfulVotes,
		review.CreatedAt,
		review.UpdatedAt,
		false,
	)

	return err
}

func (r *reviewRepository) GetByID(id uuid.UUID) (*model.Review, error) {
	var review model.Review

	query := `
		SELECT r.*, u.username, p.name as product_name
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN products p ON r.product_id = p.id
		WHERE r.id = $1 and r.is_deleted = FALSE
	`

	err := r.db.Get(&review, query, id)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *reviewRepository) GetByProductID(productID uuid.UUID) ([]model.Review, error) {
	var reviews []model.Review

	query := `
		SELECT r.*, u.username
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1 AND r.is_deleted = FALSE
		ORDER BY r.created_at DESC
	`

	err := r.db.Select(&reviews, query, productID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *reviewRepository) GetByUserID(userID uuid.UUID) ([]model.Review, error) {
	var reviews []model.Review

	query := `
		SELECT r.*, p.name as product_name
		FROM reviews r
		JOIN products p ON r.product_id = p.id
		WHERE r.user_id = $1 AND r.is_deleted = FALSE
		ORDER BY r.created_at DESC
	`

	err := r.db.Select(&reviews, query, userID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *reviewRepository) Update(review *model.Review) error {
	review.UpdatedAt = time.Now()

	query := `
		UPDATE reviews
		SET rating = $1, title = $2, comment = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(
		query,
		review.Rating,
		review.Title,
		review.Comment,
		review.UpdatedAt,
		review.ID,
	)

	return err
}

func (r *reviewRepository) Delete(id uuid.UUID) error {
	query := `UPDATE reviews SET is_deleted = TRUE WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *reviewRepository) List(offset, limit int) ([]model.Review, int, error) {
	var reviews []model.Review
	var total int

	// Get total count
	countQuery := `SELECT COUNT(*) FROM reviews`
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// Get reviews
	listQuery := `
		SELECT r.*, u.username, p.name as product_name
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN products p ON r.product_id = p.id
		ORDER BY r.created_at DESC
		LIMIT $1 OFFSET $2
	`

	err = r.db.Select(&reviews, listQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
