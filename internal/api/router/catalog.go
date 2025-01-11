package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/handler"
	"github.com/may20xx/booking/internal/utils"
)

type catalogRouter struct {
	validate *validator.Validate
	service  handler.CatalogService
}

func newCatalogRoutes() *catalogRouter {
	return &catalogRouter{
		validate: validator.New(),
		service:  handler.NewCatalogService(),
	}
}

func (r *catalogRouter) findAll(c *fiber.Ctx) error {

	page := c.Query("page")
	limit := c.Query("limit")

	result, err := r.service.FindAll(page, limit)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (r *catalogRouter) findDetail(c *fiber.Ctx) error {

	id := c.Params("id")

	result, err := r.service.FindById(id)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (r *catalogRouter) save(c *fiber.Ctx) error {
	req := new(dto.CatalogRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid request body!"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	result, err := r.service.Insert(req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (r *catalogRouter) update(c *fiber.Ctx) error {
	id := c.Params("id")

	req := new(dto.CatalogRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid request body!"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	result, err := r.service.Update(id, req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (r *catalogRouter) delete(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := r.service.Remove(id)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func CatalogRouter(r fiber.Router) {
	routes := newCatalogRoutes()

	r.Get("/catalogs", routes.findAll)
	r.Post("/catalogs", routes.save)
	r.Get("/catalogs/:id", routes.findDetail)
	r.Put("/catalogs/:id", routes.update)
	r.Delete("/catalogs/:id", routes.delete)
}
