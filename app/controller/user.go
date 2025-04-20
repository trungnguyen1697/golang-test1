package controller

import (
	"golang-test1/app/dto"
	"golang-test1/app/model"
	repo "golang-test1/app/repository"
	"golang-test1/pkg/validator"
	"golang-test1/platform/database"

	"github.com/gofiber/fiber/v2"
)

// CreateUser func for creates a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param register body model.RegisterUser true "Create new user"
// @Failure 400,401,409,500 {object} ErrorResponse "Error"
// @Success 200 {object} dto.User "Ok"
// @Router /api/v1/register [post]
func CreateUser(c *fiber.Ctx) error {
	// Create new Book struct
	user := &model.RegisterUser{}

	if err := c.BodyParser(user); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := validator.NewValidator()
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	// check user already exists
	exists, err := userRepo.Exists(user.UserName, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"msg": "user with this username or email already exists",
		})
	}

	user.Password, _ = GeneratePasswordHash([]byte(user.Password))
	if err := userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbUser, err := userRepo.GetByUsername(user.UserName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(dbUser),
	})
}

// GetUsers func gets all exists user.
// @Description Get all exists user.
// @Summary get all exists user
// @Tags User
// @Accept json
// @Produce json
// @Param page query integer false "Page no"
// @Param page_size query integer false "records per page"
// @Success 200 {array} dto.User
// @Failure 400,401,403 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func GetUsers(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)
	userRepo := repo.NewUserRepo(database.GetDB())
	users, err := userRepo.All(pageSize, uint(pageSize*(pageNo-1)))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "users were not found",
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(users),
		"users":     dto.ToUsers(users),
	})
}

// GetUser func gets a user.
// @Description a user.
// @Summary get a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)"  // Updated description
// @Success 200 {object} dto.User
// @Failure 400,401,403,404 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userRepo := repo.NewUserRepo(database.GetDB())
	user, err := userRepo.Get(idStr)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(user),
	})
}

// UpdateUser func update a user.
// @Description first_name, last_name, is_active, is_admin only
// @Summary update a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param updateuser body model.UpdateUser true "Update a user"
// @Success 200 {object} dto.User
// @Failure 400,401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userRepo := repo.NewUserRepo(database.GetDB())
	_, err := userRepo.Get(idStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	user := &model.UpdateUser{}
	if err := c.BodyParser(user); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := validator.NewValidator()
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := userRepo.Update(idStr, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbUser, err := userRepo.Get(idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(dbUser),
	})
}

// DeleteUser func delete a user.
// @Description delete user
// @Summary delete a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} interface{} "Ok"
// @Failure 401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userRepo := repo.NewUserRepo(database.GetDB())
	_, err := userRepo.Get(idStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	err = userRepo.Delete(idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{})
}

// ChangeUserPassword func changes a user's password.
// @Description Change user password with verification of current password
// @Summary change user password
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param changePassword body model.ChangePassword true "Password change data"
// @Success 200 {object} interface{} "Ok"
// @Failure 400,401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/users/{id}/change-password [post]
func ChangeUserPassword(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userRepo := repo.NewUserRepo(database.GetDB())

	// Check if user exists
	_, err := userRepo.Get(idStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user was not found",
		})
	}

	// Parse the request body
	passwordData := &model.ChangePassword{}
	if err := c.BodyParser(passwordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Validate the input
	validate := validator.NewValidator()
	if err := validate.Struct(passwordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	// Get current password hash from database
	currentPasswordHash, err := userRepo.GetPassword(idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to verify current password",
		})
	}

	// Verify current password
	if !IsValidPassword([]byte(currentPasswordHash), []byte(passwordData.CurrentPassword)) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "current password is incorrect",
		})
	}

	// Hash the new password
	newPasswordHash, err := GeneratePasswordHash([]byte(passwordData.NewPassword))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to hash new password",
		})
	}

	// Update the password in the database
	if err := userRepo.ChangePassword(idStr, newPasswordHash); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "password changed successfully",
	})
}
