package interceptor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/utils"
)

func RouteNotMatch() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SendStatus(404)

		return c.JSON(utils.NewAppError(fiber.StatusNotFound, "Not found route!"))
	}
}
