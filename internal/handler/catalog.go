package handler

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/domain"
	"github.com/may20xx/booking/internal/storage"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/log"
)

type CatalogService interface {
	FindAll(page string, limit string) (*utils.Pagination, *utils.AppError)
	FindById(id string) (*utils.Response, *utils.AppError)
	Insert(catalog *dto.CatalogRequest) (*utils.Response, *utils.AppError)
	Update(id string, catalog *dto.CatalogRequest) (*utils.Response, *utils.AppError)
	Remove(id string) (*utils.Response, *utils.AppError)
}

type catalogService struct {
	catalogRepo storage.CatalogRepository
}

func NewCatalogService() *catalogService {
	ctx := context.Background()

	db, err := database.GetDatabase(ctx)
	if err != nil {
		log.Msg.Panic("error getting database connection %s", err)
	}

	return &catalogService{catalogRepo: storage.NewCatalogRepository(db, ctx)}
}

func (s *catalogService) FindAll(page string, limit string) (*utils.Pagination, *utils.AppError) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}

	res, total, err := s.catalogRepo.FindAll(pageInt, limitInt)
	if err != nil {
		return nil, utils.NewAppError(500, "Internal server error")
	}

	return utils.NewPaginationResponse(total, pageInt, limitInt, res), nil
}

func (s *catalogService) FindById(id string) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	res, err := s.catalogRepo.FindById(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Catalog not found!")
		}
	}

	return utils.NewResponse(200, res), nil
}

func (s *catalogService) Insert(catalog *dto.CatalogRequest) (*utils.Response, *utils.AppError) {

	newCatalog := &domain.Catalog{
		Name: catalog.Name,
	}

	res, err := s.catalogRepo.Insert(newCatalog)

	if err != nil {
		return nil, utils.NewAppError(500, "Internal server error")
	}

	return utils.NewResponse(201, res), nil
}

func (s *catalogService) Update(id string, catalog *dto.CatalogRequest) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	existingCatalog, err := s.catalogRepo.FindById(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Catalog not found!")
		}
	}

	existingCatalog.Name = catalog.Name

	res, err := s.catalogRepo.Update(existingCatalog)

	if err != nil {
		return nil, utils.NewAppError(500, "Internal server error")
	}

	return utils.NewResponse(200, res), nil
}

func (s *catalogService) Remove(id string) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	err = s.catalogRepo.Remove(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Catalog not found!")
		}
	}

	return utils.NewResponse(200, "Deleted catalog successfully!"), nil
}
