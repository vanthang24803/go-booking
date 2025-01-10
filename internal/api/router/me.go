package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/api/middleware/guard"
	"github.com/may20xx/booking/internal/handler"
	"github.com/may20xx/booking/internal/utils"
)

type MeRoutes struct {
	validate *validator.Validate
	service  handler.MeHandler
}

func NewMeRoutes() *MeRoutes {
	return &MeRoutes{
		validate: validator.New(),
		service:  handler.NewMeService(),
	}
}

func (r *MeRoutes) profile(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	result, err := r.service.GetProfile(payload)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewResponse(fiber.StatusOK, result))
}

func (r *MeRoutes) uploadAvatar(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Failed to get file"))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(500, "Failed to open file"))
	}
	defer file.Close()

	res, ext := r.service.UploadAvatar(payload, file)

	if ext != nil {
		return c.Status(ext.Code).JSON(ext)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewResponse(fiber.StatusOK, res))
}

func (r *MeRoutes) updateProfile(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)

	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	req := new(dto.UpdateProfileRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	res, err := r.service.UpdateProfile(payload, req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewResponse(fiber.StatusOK, res))
}

func MeRouter(router fiber.Router) {
	routes := NewMeRoutes()

	router.Get("/me", guard.AuthGuard(), routes.profile)
	router.Post("/me/avatar", guard.AuthGuard(), routes.uploadAvatar)
	router.Put("/me", guard.AuthGuard(), routes.updateProfile)

}
