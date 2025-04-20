package controller

import (
	"github.com/golang-jwt/jwt/v4"
	"golang-test1/app/dto"
	repo "golang-test1/app/repository"
	"golang-test1/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AddToWishlist adds a product to a user's wishlist
// @Description Add a product to the user's wishlist
// @Summary add product to wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param wishlistItem body dto.WishlistItemRequest true "Product ID to add"
// @Success 200 {object} interface{} "Success"
// @Failure 400,401,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/wishlist [post]
func AddToWishlist(c *fiber.Ctx) error {
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
	var data dto.WishlistItemRequest

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid request body",
		})
	}

	// Validate product ID
	productID, err := uuid.Parse(data.ProductID)
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

	// Add to wishlist
	wishlistRepo := repo.NewWishlistRepository(database.GetDB())
	if err := wishlistRepo.Add(userID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to add product to wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "product added to wishlist successfully",
	})
}

// RemoveFromWishlist removes a product from a user's wishlist
// @Description Remove a product from the user's wishlist
// @Summary remove product from wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} interface{} "Success"
// @Failure 400,401,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/wishlist/{product_id} [delete]
func RemoveFromWishlist(c *fiber.Ctx) error {
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

	// Get product ID from path parameter
	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Remove from wishlist
	wishlistRepo := repo.NewWishlistRepository(database.GetDB())

	// Check if product is in wishlist
	exists, err := wishlistRepo.IsInWishlist(userID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to check wishlist",
		})
	}

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found in wishlist",
		})
	}

	if err := wishlistRepo.Remove(userID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to remove product from wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "product removed from wishlist successfully",
	})
}

// GetWishlist returns all items in a user's wishlist
// @Description Get all items in the user's wishlist
// @Summary get user wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Success 200 {array} model.WishlistItem "Wishlist items"
// @Failure 401,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/wishlist [get]
func GetWishlist(c *fiber.Ctx) error {
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

	// Get wishlist items
	wishlistRepo := repo.NewWishlistRepository(database.GetDB())
	items, err := wishlistRepo.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to retrieve wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"count": len(items),
		"items": items,
	})
}

// CheckWishlistItem checks if a product is in a user's wishlist
// @Description Check if a product is in the user's wishlist
// @Summary check wishlist item
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} map[string]bool "Result with in_wishlist field"
// @Failure 400,401,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/wishlist/check/{product_id} [get]
func CheckWishlistItem(c *fiber.Ctx) error {
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

	// Get product ID from path parameter
	productIDStr := c.Params("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Check wishlist
	wishlistRepo := repo.NewWishlistRepository(database.GetDB())
	inWishlist, err := wishlistRepo.IsInWishlist(userID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to check wishlist",
		})
	}

	return c.JSON(fiber.Map{
		"in_wishlist": inWishlist,
	})
}
