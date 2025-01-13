package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/handler"
)

type bookingRouter struct {
	validate *validator.Validate
	service  handler.BookingService
}

func newBookingRouter() *bookingRouter {
	return &bookingRouter{
		validate: validator.New(),
		service:  handler.NewBookingService(),
	}
}

func (r *bookingRouter) findAll(c *fiber.Ctx) error {
	return nil
}

func BookingRouter(router fiber.Router) {
	routes := newBookingRouter()

	router.Get("/bookings", routes.findAll)
}
