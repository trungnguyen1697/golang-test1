package controller

import (
	repo "golang-test1/app/repository"
	"golang-test1/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// GetDashboardStats returns user dashboard statistics
// @Description Get dashboard statistics for current user
// @Summary get user dashboard stats
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} model.DashboardStats "Dashboard statistics"
// @Failure 401,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/stats [get]
func GetDashboardStats(c *fiber.Ctx) error {
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

	// Get dashboard stats
	wishlistRepo := repo.NewWishlistRepository(database.GetDB())
	stats, err := wishlistRepo.GetDashboardStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg":   "failed to retrieve dashboard statistics",
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}
