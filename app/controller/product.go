package controller

import (
	"golang-test1/app/dto"
	"golang-test1/app/model"
	repo "golang-test1/app/repository"
	"golang-test1/pkg/validator"
	"golang-test1/platform/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateProduct func for creating a new product.
// @Description Create a new product.
// @Summary create a new product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body model.ProductInput true "Create new product"
// @Success 201 {object} dto.Product "Created"
// @Failure 400,401,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/products [post]
func CreateProduct(c *fiber.Ctx) error {
	// Parse request body
	productInput := &model.ProductInput{}
	if err := c.BodyParser(productInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Validate input
	validate := validator.NewValidator()
	if err := validate.Struct(productInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	// Create product object
	now := time.Now()
	product := &model.Product{
		ID:            uuid.New(),
		SKU:           productInput.SKU,
		Name:          productInput.Name,
		Description:   productInput.Description,
		Price:         productInput.Price,
		StockQuantity: productInput.StockQuantity,
		Status:        productInput.Status,
		Attributes:    productInput.Attributes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Handle optional fields
	if productInput.SalePrice != nil {
		product.SalePrice = *productInput.SalePrice
	}
	if productInput.CostPrice != nil {
		product.CostPrice = *productInput.CostPrice
	}

	// Save product to database
	productRepo := repo.NewProductRepository(database.GetDB())
	if err := productRepo.Create(product, productInput.CategoryIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Get complete product with categories
	createdProduct, err := productRepo.GetByID(product.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"product": dto.ToProduct(createdProduct),
	})
}

// GetProduct func gets a single product by ID.
// @Description Get product details by ID.
// @Summary get a product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID format)"
// @Success 200 {object} dto.Product
// @Failure 400,404 {object} ErrorResponse "Error"
// @Router /api/v1/products/{id} [get]
func GetProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	productRepo := repo.NewProductRepository(database.GetDB())
	product, err := productRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	return c.JSON(fiber.Map{
		"product": dto.ToProduct(product),
	})
}

// UpdateProduct func for updating a product.
// @Description Update an existing product.
// @Summary update a product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID format)"
// @Param product body model.ProductInput true "Update product"
// @Success 200 {object} dto.Product "Ok"
// @Failure 400,401,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Parse request body
	productInput := &model.ProductInput{}
	if err := c.BodyParser(productInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Validate input
	validate := validator.NewValidator()
	if err := validate.Struct(productInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	// Check if product exists
	productRepo := repo.NewProductRepository(database.GetDB())
	existingProduct, err := productRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	// Update product
	existingProduct.SKU = productInput.SKU
	existingProduct.Name = productInput.Name
	existingProduct.Description = productInput.Description
	existingProduct.Price = productInput.Price
	existingProduct.StockQuantity = productInput.StockQuantity
	existingProduct.Status = productInput.Status
	existingProduct.Attributes = productInput.Attributes
	existingProduct.UpdatedAt = time.Now()

	// Handle optional fields
	if productInput.SalePrice != nil {
		existingProduct.SalePrice = *productInput.SalePrice
	}
	if productInput.CostPrice != nil {
		existingProduct.CostPrice = *productInput.CostPrice
	}

	// Save updated product
	if err := productRepo.Update(existingProduct, productInput.CategoryIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Get updated product with categories
	updatedProduct, err := productRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"product": dto.ToProduct(updatedProduct),
	})
}

// DeleteProduct func for deleting a product.
// @Description Delete a product.
// @Summary delete a product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID format)"
// @Success 200 {object} interface{} "Ok"
// @Failure 400,401,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/products/{id} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Check if product exists
	productRepo := repo.NewProductRepository(database.GetDB())
	_, err = productRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	// Delete product
	if err := productRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "product successfully deleted",
	})
}

// GetProducts func for listing products with filtering and pagination.
// @Description List all products with optional filtering and pagination.
// @Summary list products
// @Tags Product
// @Accept json
// @Produce json
// @Param page query integer false "Page number"
// @Param page_size query integer false "Page size"
// @Param search query string false "Search term for name, SKU or description"
// @Param category_id query string false "Filter by category ID (UUID format)"
// @Param status query string false "Filter by status (active, inactive, out_of_stock)"
// @Param min_price query number false "Filter by minimum price"
// @Param max_price query number false "Filter by maximum price"
// @Param min_stock query integer false "Filter by minimum stock quantity"
// @Param max_stock query integer false "Filter by maximum stock quantity"
// @Param created_after query string false "Filter by creation date (RFC3339 format)"
// @Param created_before query string false "Filter by creation date (RFC3339 format)"
// @Param sort_by query string false "Sort field (name, price, created_at, stock_quantity)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {array} dto.Product
// @Failure 400,500 {object} ErrorResponse "Error"
// @Router /api/v1/products [get]
func GetProducts(c *fiber.Ctx) error {
	// Get pagination parameters
	pageNo, pageSize := GetPagination(c)
	offset := (pageNo - 1) * pageSize

	// Get filter parameters
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Parse category ID if provided
	var categoryID *uuid.UUID
	if c.Query("category_id") != "" {
		id, err := uuid.Parse(c.Query("category_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid category ID format",
			})
		}
		categoryID = &id
	}

	// Parse price range parameters
	var minPrice, maxPrice *float64
	if c.Query("min_price") != "" {
		val, err := strconv.ParseFloat(c.Query("min_price"), 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid min_price format",
			})
		}
		minPrice = &val
	}
	if c.Query("max_price") != "" {
		val, err := strconv.ParseFloat(c.Query("max_price"), 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid max_price format",
			})
		}
		maxPrice = &val
	}

	// Parse stock quantity range parameters
	var minStock, maxStock *int
	if c.Query("min_stock") != "" {
		val, err := strconv.Atoi(c.Query("min_stock"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid min_stock format",
			})
		}
		minStock = &val
	}
	if c.Query("max_stock") != "" {
		val, err := strconv.Atoi(c.Query("max_stock"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid max_stock format",
			})
		}
		maxStock = &val
	}

	// Parse date range parameters
	var createdAfter, createdBefore *time.Time
	if c.Query("created_after") != "" {
		t, err := time.Parse(time.RFC3339, c.Query("created_after"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid created_after format, use RFC3339",
			})
		}
		createdAfter = &t
	}
	if c.Query("created_before") != "" {
		t, err := time.Parse(time.RFC3339, c.Query("created_before"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": "invalid created_before format, use RFC3339",
			})
		}
		createdBefore = &t
	}

	// Validate sort parameters
	validSortFields := map[string]bool{
		"name":           true,
		"price":          true,
		"created_at":     true,
		"stock_quantity": true,
		"sku":            true,
	}
	validSortOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}
	if !validSortOrders[sortOrder] {
		sortOrder = "desc"
	}

	// Get products from repository with enhanced filtering
	productRepo := repo.NewProductRepository(database.GetDB())
	products, total, err := productRepo.ListWithFilters(
		offset, pageSize, search, categoryID, status,
		minPrice, maxPrice, minStock, maxStock,
		createdAfter, createdBefore, sortBy, sortOrder,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"total":     total,
		"products":  dto.ToProducts(products),
	})
}

// GetProductCategories func for getting categories of a product.
// @Description Get all categories of a product.
// @Summary get product categories
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID format)"
// @Success 200 {array} dto.Category
// @Failure 400,404 {object} ErrorResponse "Error"
// @Router /api/v1/products/{id}/categories [get]
func GetProductCategories(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "invalid product ID format",
		})
	}

	// Check if product exists
	productRepo := repo.NewProductRepository(database.GetDB())
	_, err = productRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "product not found",
		})
	}

	// Get product categories
	categories, err := productRepo.GetCategories(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"categories": dto.ToCategories(categories),
	})
}
