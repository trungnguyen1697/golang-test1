package controller

import (
	"golang-test1/app/dto"
	"golang-test1/app/model"
	repo "golang-test1/app/repository"
	"golang-test1/pkg/validator"
	"golang-test1/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// CatErrorResponse represents an error response.
type CatErrorResponse struct {
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// CategoryResponse represents a successful category response.
type CategoryResponse struct {
	Category dto.Category `json:"category"`
}

// CategoriesResponse represents a successful categories list response.
type CategoriesResponse struct {
	Categories []dto.Category `json:"categories"`
}

// ProductCountResponse represents a product count response.
type ProductCountResponse struct {
	ProductCount int `json:"product_count"`
}

// SuccessResponse represents a simple success message response.
type SuccessResponse struct {
	Message string `json:"msg"`
}

// CreateCategory func for creating a new category.
// @Description Create a new category.
// @Summary create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body model.CategoryInput true "Create new category"
// @Success 201 {object} CategoryResponse
// @Failure 400,401,500 {object} CatErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/categories [post]
func CreateCategory(c *fiber.Ctx) error {
	// Create new CategoryInput struct
	input := &model.CategoryInput{}

	// Parse request body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
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

	// Create category entity
	now := time.Now()
	category := &model.Category{
		ID:           uuid.New(),
		Name:         input.Name,
		Slug:         slug.Make(input.Name),
		Description:  input.Description,
		ParentID:     input.ParentID,
		IsActive:     input.IsActive,
		DisplayOrder: input.DisplayOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Create category in database
	if err := categoryRepo.Create(category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return created category
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"category": dto.ToCategory(category),
	})
}

// GetCategory func gets a category by ID.
// @Description Get a category by ID.
// @Summary get a category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Success 200 {object} CategoryResponse
// @Failure 400,404 {object} CatErrorResponse "Error"
// @Router /api/v1/categories/{id} [get]
func GetCategory(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid category ID format",
		})
	}

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Get category from database
	category, err := categoryRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "category not found",
		})
	}

	// Return category
	return c.JSON(fiber.Map{
		"category": dto.ToCategory(category),
	})
}

// GetCategoryBySlug func gets a category by slug.
// @Description Get a category by slug.
// @Summary get a category by slug
// @Tags Category
// @Accept json
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} CategoryResponse
// @Failure 404 {object} CatErrorResponse "Error"
// @Router /api/v1/categories/slug/{slug} [get]
func GetCategoryBySlug(c *fiber.Ctx) error {
	// Parse slug from URL parameter
	slugStr := c.Params("slug")

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Get category from database
	category, err := categoryRepo.GetBySlug(slugStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "category not found",
		})
	}

	// Return category
	return c.JSON(fiber.Map{
		"category": dto.ToCategory(category),
	})
}

// UpdateCategory func updates a category.
// @Description Update a category.
// @Summary update a category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Param category body model.CategoryInput true "Update category"
// @Success 200 {object} CategoryResponse
// @Failure 400,401,404,500 {object} CatErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/categories/{id} [put]
func UpdateCategory(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid category ID format",
		})
	}

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Check if category exists
	existingCategory, err := categoryRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "category not found",
		})
	}

	// Parse request body
	input := &model.CategoryInput{}
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
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

	// Update category properties
	existingCategory.Name = input.Name
	existingCategory.Slug = slug.Make(input.Name)
	existingCategory.Description = input.Description
	existingCategory.ParentID = input.ParentID
	existingCategory.IsActive = input.IsActive
	existingCategory.DisplayOrder = input.DisplayOrder

	// Update category in database
	if err := categoryRepo.Update(existingCategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Get updated category
	updatedCategory, err := categoryRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return updated category
	return c.JSON(fiber.Map{
		"category": dto.ToCategory(updatedCategory),
	})
}

// DeleteCategory func deletes a category.
// @Description Delete a category.
// @Summary delete a category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Success 200 {object} SuccessResponse "success message"
// @Failure 400,401,404,409,500 {object} CatErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/categories/{id} [delete]
func DeleteCategory(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid category ID format",
		})
	}

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Check if category exists
	_, err = categoryRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "category not found",
		})
	}

	// Delete category
	err = categoryRepo.Delete(id)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "category has products" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"msg": "category cannot be deleted because it has products",
			})
		} else if err.Error() == "category has child categories" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"msg": "category cannot be deleted because it has child categories",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return success message
	return c.JSON(fiber.Map{
		"msg": "category deleted successfully",
	})
}

// ListCategories func gets all categories.
// @Description Get all categories.
// @Summary get all categories
// @Tags Category
// @Accept json
// @Produce json
// @Success 200 {object} CategoriesResponse
// @Failure 500 {object} CatErrorResponse "Error"
// @Router /api/v1/categories [get]
func ListCategories(c *fiber.Ctx) error {
	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Get all categories
	categories, err := categoryRepo.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return categories
	return c.JSON(fiber.Map{
		"categories": dto.ToCategories(categories),
	})
}

// GetCategoryProductCount func gets the count of products in a category.
// @Description Get the count of products in a category.
// @Summary get product count for a category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Success 200 {object} ProductCountResponse "product count"
// @Failure 400,404,500 {object} CatErrorResponse "Error"
// @Router /api/v1/categories/{id}/product-count [get]
func GetCategoryProductCount(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid category ID format",
		})
	}

	// Get repository
	categoryRepo := repo.NewCategoryRepository(database.GetDB())

	// Check if category exists
	_, err = categoryRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "category not found",
		})
	}

	// Get product count
	count, err := categoryRepo.GetProductCount(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return product count
	return c.JSON(fiber.Map{
		"product_count": count,
	})
}
