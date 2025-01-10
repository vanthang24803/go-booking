package interceptor

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/pkg/log"
)

func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		err := c.Next()

		duration := time.Since(startTime)

		statusCode := c.Response().StatusCode()

		log.Msg.Infof("%s %s - %d - %v", c.Method(), c.Path(), statusCode, duration)

		if err != nil {
			log.Msg.With("error", err.Error()).Error("Request failed!")
		}

		return err
	}
}
