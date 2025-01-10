package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/handler"
	"github.com/may20xx/booking/internal/utils"
)

type AuthRoutes struct {
	validate *validator.Validate
	service  handler.AuthHandler
}

func NewAuthRoutes() *AuthRoutes {
	return &AuthRoutes{
		validate: validator.New(),
		service:  handler.NewAuthService(),
	}
}

func (r *AuthRoutes) register(c *fiber.Ctx) error {
	req := new(dto.RegisterRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	result, err := r.service.RegisterHandler(req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(utils.NewResponse(fiber.StatusCreated, result))
}

func (r *AuthRoutes) login(c *fiber.Ctx) error {
	req := new(dto.LoginRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	result, err := r.service.LoginHandler(req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewResponse(fiber.StatusOK, result))
}

func AuthRouter(router fiber.Router) {
	routes := NewAuthRoutes()

	router.Post("/auth/register", routes.register)
	router.Post("/auth/login", routes.login)

}
