package middleware

import (
	"errors"
	"golang-test1/pkg/config"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3" // Update to v3 which uses golang-jwt/jwt/v4
	"github.com/golang-jwt/jwt/v4"      // Replace form3tech-oss/jwt-go with this
)

// JWTProtected func for specify route group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	jwtwareConfig := jwtware.Config{
		SigningKey:     []byte(config.AppCfg().JWTSecretKey),
		ContextKey:     "user", // used in private route
		ErrorHandler:   jwtError,
		SuccessHandler: verifyTokenExpiration,
		SigningMethod:  "HS256",
	}

	return jwtware.New(jwtwareConfig)
}

func verifyTokenExpiration(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	expires := int64(claims["exp"].(float64))
	if time.Now().Unix() > expires {
		return jwtError(c, errors.New("token expired"))
	}

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"msg": err.Error(),
	})
}
