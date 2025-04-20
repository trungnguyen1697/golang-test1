package controller

import (
	"fmt"
	"golang-test1/app/model"
	repo "golang-test1/app/repository"
	"golang-test1/pkg/config"
	"golang-test1/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login method for user authentication.
// @Description Authenticate user and return an access token.
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body model.LoginUser true "Login credentials"
// @Failure 400,404,401,500 {object} ErrorResponse "Error"
// @Success 200 {object} TokenResponse "Ok"
// @Router /api/v1/login [post]
func Login(c *fiber.Ctx) error {
	login := &model.LoginUser{}
	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	userRepo := repo.NewUserRepo(database.GetDB())
	user, err := userRepo.GetByUsername(login.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "username not found",
		})
	}

	isValid := IsValidPassword([]byte(user.PasswordHash), []byte(login.Password))
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "password is wrong",
		})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "user not active anymore.",
		})
	}

	// Generate a new Access token.
	token, err := GenerateNewAccessToken(user.ID, user.Role)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg":          fmt.Sprintf("Token will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"access_token": token,
	})
}

func GenerateNewAccessToken(userID uuid.UUID, role string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID.String() // Explicitly convert UUID to string
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.AppCfg().JWTSecretExpireMinutesCount)).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.AppCfg().JWTSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GeneratePasswordHash(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func IsValidPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false
	}

	return true
}
