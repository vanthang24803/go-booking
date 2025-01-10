package guard

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/utils"
)

func AuthGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Missing authorization header"))
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Wrong authorization header format"))
		}

		token := bearerToken[1]

		payload, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Invalid token"))
		}

		c.Locals("user", payload)

		return c.Next()
	}
}
