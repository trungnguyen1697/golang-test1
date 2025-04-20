package route

import (
	"golang-test1/app/controller"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public route.
func PublicRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/api/v1")

	route.Post("/login", controller.Login)
	route.Post("/register", controller.CreateUser)

	// Category route group - Public routes for viewing categories
	categoryRoute := a.Group("/api/v1/categories")
	categoryRoute.Get("/", controller.ListCategories)                           // Get all categories
	categoryRoute.Get("/:id", controller.GetCategory)                           // Get a category by ID
	categoryRoute.Get("/slug/:slug", controller.GetCategoryBySlug)              // Get a category by slug
	categoryRoute.Get("/:id/product-count", controller.GetCategoryProductCount) // Get product count for a category

	// Public product routes - accessible without authentication
	productPublicRoute := a.Group("/api/v1/products")
	productPublicRoute.Get("/", controller.GetProducts)                        // List all products
	productPublicRoute.Get("/:id", controller.GetProduct)                      // Get a product by ID
	productPublicRoute.Get("/:id/categories", controller.GetProductCategories) // Get product categories
}
