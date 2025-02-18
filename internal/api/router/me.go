package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/api/middleware/guard"
	"github.com/may20xx/booking/internal/handler"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/log"
)

type meRouter struct {
	validate *validator.Validate
	service  handler.MeService
}

func newMeRouter() *meRouter {
	return &meRouter{
		validate: validator.New(),
		service:  handler.NewMeService(),
	}
}

func (r *meRouter) profile(c *fiber.Ctx) error {
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

func (r *meRouter) uploadAvatar(c *fiber.Ctx) error {
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

	log.Msg.Debug(file)

	res, ext := r.service.UploadAvatar(payload, file)

	if ext != nil {
		log.Msg.Error(ext.Error())
		return c.Status(ext.Code).JSON(ext)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewResponse(fiber.StatusOK, res))
}

func (r *meRouter) updateProfile(c *fiber.Ctx) error {
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

func (r *meRouter) logout(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)

	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	err := r.service.Logout(payload)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewAppError(fiber.StatusOK, "Logout successfully"))
}

func MeRouter(router fiber.Router) {
	routes := newMeRouter()

	router.Get("/me", guard.AuthGuard(), routes.profile)
	router.Post("/me/avatar", guard.AuthGuard(), routes.uploadAvatar)
	router.Put("/me", guard.AuthGuard(), routes.updateProfile)
	router.Post("/me/logout", guard.AuthGuard(), routes.logout)

}
