package helper

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return SendUnauthorizedResponse(c)
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return SendErrorResponse(c, fiber.StatusUnauthorized, "Format token tidak valid", nil)
		}

		claims, err := ValidateAccessToken(tokenParts[1], secret)
		if err != nil {
			return SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}

func GetUserIDFromToken(c *fiber.Ctx) (uint, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, errors.New("user ID tidak ditemukan dalam token")
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, errors.New("user ID tidak valid")
}
