package controller

import (
	"time"

	"golang-test1/app/model"
	repo "golang-test1/app/repository"
	"golang-test1/pkg/validator"
	"golang-test1/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// CreateReview creates a new product review
// @Description Add a new review for a product
// @Summary create product review
// @Tags Review
// @Accept json
// @Produce json
// @Param reviewInput body model.ReviewInput true "Review details"
// @Success 200 {object} model.Review "Success"
// @Failure 400,401,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/reviews [post]
func CreateReview(c *fiber.Ctx) error {
	// Get user ID from JWT context
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDValue := claims["user_id"]

	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, user ID not found",
		})
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID format",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID",
		})
	}

	// Parse request body
	var input model.ReviewInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid request body",
		})
	}

	// Validate input
	validate := validator.NewValidator()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	// Verify product exists
	productRepo := repo.NewProductRepository(database.GetDB())
	if _, err := productRepo.GetByID(input.ProductID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	// Create review object
	now := time.Now()
	review := &model.Review{
		ID:                 uuid.New(),
		ProductID:          input.ProductID,
		UserID:             userID,
		Rating:             input.Rating,
		Title:              input.Title,
		Comment:            input.Comment,
		IsVerifiedPurchase: false, // This could be determined by checking order history
		HelpfulVotes:       0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Save to database
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	if err := reviewRepo.Create(review); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to create review",
		})
	}

	// Get the complete review with username and product name
	createdReview, err := reviewRepo.GetByID(review.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "review created but failed to retrieve details",
		})
	}

	return c.JSON(fiber.Map{
		"msg":    "review created successfully",
		"review": createdReview,
	})
}

// GetReview retrieves a review by ID
// @Description Get a review by its ID
// @Summary get review by ID
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} model.Review "Review details"
// @Failure 400,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/reviews/{id} [get]
func GetReview(c *fiber.Ctx) error {
	// Parse review ID from URL parameter
	reviewIDStr := c.Params("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid review ID format",
		})
	}

	// Get review from database
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	review, err := reviewRepo.GetByID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "review not found",
		})
	}

	return c.JSON(fiber.Map{
		"review": review,
	})
}

// GetProductReviews retrieves all reviews for a product
// @Description Get all reviews for a specific product
// @Summary get product reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {array} model.Review "Product reviews"
// @Failure 400,404,500 {object} ErrorResponse "Error"
// @Router /api/v1/products/{product_id}/reviews [get]
func GetProductReviews(c *fiber.Ctx) error {
	// Parse product ID from URL parameter
	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Verify product exists
	productRepo := repo.NewProductRepository(database.GetDB())
	if _, err := productRepo.GetByID(productID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	// Get reviews for the product
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	reviews, err := reviewRepo.GetByProductID(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to retrieve product reviews",
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(reviews),
		"reviews": reviews,
	})
}

// GetUserReviews retrieves all reviews by a user
// @Description Get all reviews written by the current user
// @Summary get user reviews
// @Tags Review
// @Accept json
// @Produce json
// @Success 200 {array} model.Review "User reviews"
// @Failure 401,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/reviews/my-reviews [get]
func GetUserReviews(c *fiber.Ctx) error {
	// Get user ID from JWT context
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDValue := claims["user_id"]
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, user ID not found",
		})
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID format",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID",
		})
	}

	// Get reviews by the user
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	reviews, err := reviewRepo.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to retrieve user reviews",
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(reviews),
		"reviews": reviews,
	})
}

// UpdateReview updates an existing review
// @Description Update an existing review
// @Summary update review
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Param reviewUpdate body model.ReviewInput true "Updated review details"
// @Success 200 {object} model.Review "Updated review"
// @Failure 400,401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/reviews/{id} [put]
func UpdateReview(c *fiber.Ctx) error {
	// Get user ID from JWT context
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDValue := claims["user_id"]
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, user ID not found",
		})
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID format",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID",
		})
	}

	// Parse review ID from URL parameter
	reviewIDStr := c.Params("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid review ID format",
		})
	}

	// Get the existing review
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	existingReview, err := reviewRepo.GetByID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "review not found",
		})
	}

	// Verify the user is the owner of the review
	if existingReview.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "you can only update your own reviews",
		})
	}

	// Parse request body
	var input model.ReviewInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid request body",
		})
	}

	// Validate input
	validate := validator.NewValidator()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	// Update review fields
	existingReview.Rating = input.Rating
	existingReview.Title = input.Title
	existingReview.Comment = input.Comment
	existingReview.UpdatedAt = time.Now()

	// Save to database
	if err := reviewRepo.Update(existingReview); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to update review",
		})
	}

	// Get the updated review
	updatedReview, err := reviewRepo.GetByID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "review updated but failed to retrieve details",
		})
	}

	return c.JSON(fiber.Map{
		"msg":    "review updated successfully",
		"review": updatedReview,
	})
}

// DeleteReview deletes a review
// @Description Delete a review
// @Summary delete review
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} interface{} "Success"
// @Failure 400,401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/reviews/{id} [delete]
func DeleteReview(c *fiber.Ctx) error {
	// Get user ID from JWT context
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDValue := claims["user_id"]
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, user ID not found",
		})
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID format",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized, invalid user ID",
		})
	}

	// Parse review ID from URL parameter
	reviewIDStr := c.Params("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid review ID format",
		})
	}

	// Get the existing review
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	existingReview, err := reviewRepo.GetByID(reviewID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "review not found",
		})
	}

	// Verify the user is the owner of the review (or admin)
	isAdmin := claims["is_admin"] != nil && claims["is_admin"].(bool)
	if existingReview.UserID != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "you can only delete your own reviews",
		})
	}

	// Delete the review
	if err := reviewRepo.Delete(reviewID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to delete review",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "review deleted successfully",
	})
}

// ListReviews gets all reviews with pagination
// @Description Get all reviews with pagination
// @Summary list all reviews
// @Tags Review
// @Accept json
// @Produce json
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Success 200 {array} model.Review "Reviews list"
// @Failure 500 {object} ErrorResponse "Error"
// @Router /api/v1/reviews [get]
func ListReviews(c *fiber.Ctx) error {
	// Get pagination parameters
	pageNo, pageSize := GetPagination(c)
	offset := (pageNo - 1) * pageSize

	// Get reviews from database
	reviewRepo := repo.NewReviewRepository(database.GetDB())
	reviews, total, err := reviewRepo.List(offset, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to retrieve reviews",
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"total":     total,
		"count":     len(reviews),
		"reviews":   reviews,
	})
}
