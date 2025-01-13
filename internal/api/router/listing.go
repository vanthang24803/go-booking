package router

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/api/middleware/guard"
	"github.com/may20xx/booking/internal/handler"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/log"
)

type listingRouter struct {
	validate *validator.Validate
	service  handler.ListingService
}

func newListingRouter() *listingRouter {
	return &listingRouter{
		validate: validator.New(),
		service:  handler.NewListingService(),
	}
}

func (l *listingRouter) save(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	form, err := c.MultipartForm()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid form data!"))
	}

	photos := form.File["photos"]

	img, err := validationPhotos(photos)

	if err != nil {
		log.Msg.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	request, err := validationFormData(form)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	res, err := l.service.Save(payload, request, img)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewResponse(201, res))
}

func (r *listingRouter) findDetail(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := r.service.FindDetail(id)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(result)
}

func (r *listingRouter) update(c *fiber.Ctx) error {
	id := c.Params("id")

	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	req := new(dto.ListingRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid request body!"))
	}

	if err := r.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	res, err := r.service.Update(id, req)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}
	return c.JSON(res)
}

func (r *listingRouter) remove(c *fiber.Ctx) error {
	id := c.Params("id")

	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	res, err := r.service.Remove(id)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(res)
}

func (r *listingRouter) findAll(c *fiber.Ctx) error {
	page := c.Query("page")
	limit := c.Query("limit")

	result, err := r.service.FindAll(page, limit)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(result)
}

func (r *listingRouter) searchByLocation(c *fiber.Ctx) error {
	query := c.Query("s")
	page := c.Query("page")
	limit := c.Query("limit")

	res, err := r.service.Search(page, limit, query)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(res)
}

func ListingRouter(router fiber.Router) {
	routes := newListingRouter()

	router.Get("/listings", newListingRouter().findAll)
	router.Get("/search", newListingRouter().searchByLocation)
	router.Get("/listings/:id", newListingRouter().findDetail)
	router.Post("/listings", guard.AuthGuard(), routes.save)
	router.Put("/listings/:id", guard.AuthGuard(), newListingRouter().update)
	router.Delete("/listings/:id", guard.AuthGuard(), newListingRouter().remove)
}

// Validation
const (
	maxFileSize = 5 * 1024 * 1024
)

func validationPhotos(files []*multipart.FileHeader) ([]multipart.File, error) {
	if len(files) == 0 {
		return nil, utils.NewAppError(400, "Photos are required!")
	}

	var validFiles []multipart.File

	for _, file := range files {
		if file.Size > maxFileSize {
			return nil, utils.NewAppError(400, fmt.Sprintf("File %s exceeds the size limit of 5MB", file.Filename))
		}

		f, err := file.Open()

		if err != nil {
			return nil, utils.NewAppError(400, fmt.Sprintf("Could not open file %s: %v", file.Filename, err))
		}

		defer f.Close()

		validFiles = append(validFiles, f)
	}

	return validFiles, nil
}

func validationFormData(form *multipart.Form) (*dto.ListingRequest, error) {
	getFirstValue := func(key string) (string, bool) {
		if values, exists := form.Value[key]; exists && len(values) > 0 {
			return values[0], true
		}
		return "", false
	}

	getIntValue := func(key string) (int, error) {
		value, ok := getFirstValue(key)
		if !ok {
			return 0, fmt.Errorf("%s is required", key)
		}
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("%s must be an integer", key)
		}
		return intValue, nil
	}

	getFloatValue := func(key string) (float64, error) {
		value, ok := getFirstValue(key)
		if !ok {
			return 0, fmt.Errorf("%s is required", key)
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a number", key)
		}
		return floatValue, nil
	}

	convertStrToInt := func(arr []string) ([]int, error) {
		seen := make(map[int]bool)
		var result []int

		for i, str := range arr {
			num, err := strconv.Atoi(str)
			if err != nil {
				return nil, fmt.Errorf("error converting string to int at index %d: %w", i, err)
			}
			if !seen[num] {
				seen[num] = true
				result = append(result, num)
			}
		}
		return result, nil
	}

	catalogs, err := convertStrToInt(form.Value["catalogs"])

	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, "Catalogs is required!")
	}

	title, ok := getFirstValue("title")
	if !ok {
		return &dto.ListingRequest{}, utils.NewAppError(400, "Title is required!")
	}

	description, ok := getFirstValue("description")
	if !ok {
		return &dto.ListingRequest{}, utils.NewAppError(400, "Description is required!")
	}

	location, ok := getFirstValue("location")
	if !ok {
		return &dto.ListingRequest{}, utils.NewAppError(400, "Location is required!")
	}

	guests, err := getIntValue("guests")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	beds, err := getIntValue("beds")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	baths, err := getIntValue("baths")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	price, err := getFloatValue("price")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	cleaningFee, err := getFloatValue("cleaning_fee")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	serviceFee, err := getFloatValue("service_fee")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	taxes, err := getFloatValue("taxes")
	if err != nil {
		return &dto.ListingRequest{}, utils.NewAppError(400, err.Error())
	}

	request := &dto.ListingRequest{
		Title:       title,
		Description: description,
		Location:    location,
		Guests:      guests,
		Beds:        beds,
		Baths:       baths,
		Price:       price,
		CleaningFee: cleaningFee,
		ServiceFee:  serviceFee,
		Taxes:       taxes,
		Catalogs:    catalogs,
	}

	return request, nil
}
