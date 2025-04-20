package route

import (
	"golang-test1/app/controller"
	"golang-test1/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private route.
func PrivateRoutes(a *fiber.App) {
	// Admin route group
	adminRoute := a.Group("/api/v1/users", middleware.JWTProtected(), middleware.IsAdmin)
	// User
	adminRoute.Get("/", controller.GetUsers)
	adminRoute.Get("/:id", controller.GetUser)
	adminRoute.Put("/:id", controller.UpdateUser)
	adminRoute.Delete("/:id", controller.DeleteUser)

	// User route group - routes accessible to authenticated users
	userRoute := a.Group("/api/v1/users", middleware.JWTProtected())
	// Change password endpoint
	userRoute.Post("/:id/change-password", controller.ChangeUserPassword)

	// Category route group - Admin routes for managing categories
	categoryAdminRoute := a.Group("/api/v1/categories", middleware.JWTProtected(), middleware.IsAdmin)
	categoryAdminRoute.Post("/", controller.CreateCategory)      // Create a new category
	categoryAdminRoute.Put("/:id", controller.UpdateCategory)    // Update a category
	categoryAdminRoute.Delete("/:id", controller.DeleteCategory) // Delete a category

	// Product routes
	// Admin product routes - require admin privileges
	productAdminRoute := a.Group("/api/v1/products", middleware.JWTProtected(), middleware.IsAdmin)
	productAdminRoute.Post("/", controller.CreateProduct)      // Create a new product
	productAdminRoute.Put("/:id", controller.UpdateProduct)    // Update a product
	productAdminRoute.Delete("/:id", controller.DeleteProduct) // Delete a product

	// Public product routes - accessible to all authenticated users
	productRoute := a.Group("/api/v1/products", middleware.JWTProtected())
	productRoute.Get("/", controller.GetProducts)                        // List all products
	productRoute.Get("/:id", controller.GetProduct)                      // Get a product by ID
	productRoute.Get("/:id/categories", controller.GetProductCategories) // Get product categories

	// Wishlist routes - accessible to all authenticated users
	wishlistRoute := a.Group("/api/v1/wishlist", middleware.JWTProtected())
	wishlistRoute.Post("/", controller.AddToWishlist)                     // Add item to wishlist
	wishlistRoute.Delete("/:product_id", controller.RemoveFromWishlist)   // Remove item from wishlist
	wishlistRoute.Get("/", controller.GetWishlist)                        // Get user's wishlist
	wishlistRoute.Get("/check/:product_id", controller.CheckWishlistItem) // Check if item is in wishlist

	// Review routes - accessible to all authenticated users
	reviewRoute := a.Group("/api/v1/reviews", middleware.JWTProtected())
	reviewRoute.Post("/", controller.CreateReview)            // Create a new review
	reviewRoute.Get("/", controller.ListReviews)              // List all reviews with pagination
	reviewRoute.Get("/my-reviews", controller.GetUserReviews) // Get all reviews by current user
	reviewRoute.Get("/:id", controller.GetReview)             // Get a review by ID

	reviewRoute.Put("/:id", controller.UpdateReview)    // Update a review
	reviewRoute.Delete("/:id", controller.DeleteReview) // Delete a review

	// Product review routes
	productRoute.Get("/:product_id/reviews", controller.GetProductReviews) // Get all reviews for a product

	dashboardRoutes := a.Group("/api/v1/dashboard", middleware.JWTProtected())
	dashboardRoutes.Get("/stats", controller.GetDashboardStats)
}
