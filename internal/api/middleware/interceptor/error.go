package interceptor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/utils"
)

func Error() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			if appError, ok := err.(*utils.AppError); ok {
				return c.Status(appError.Code).JSON(appError)
			}

			return c.Status(fiber.StatusInternalServerError).JSON(
				utils.NewAppError(fiber.StatusInternalServerError, "Internal server error"),
			)
		}

		return nil
	}
}
