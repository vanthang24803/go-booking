package handler

import (
	"context"
	"database/sql"
	"mime/multipart"
	"strconv"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/domain"
	"github.com/may20xx/booking/internal/storage"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/cloudinary"
	"github.com/may20xx/booking/pkg/log"
)

type ListingService interface {
	FindAll(page string, limit string) (*utils.Pagination, *utils.AppError)
	FindDetail(id string) (*utils.Response, *utils.AppError)
	Save(payload *utils.JwtPayload, req *dto.ListingRequest, files []multipart.File) (*utils.Response, error)
	Update(id string, req *dto.ListingRequest) (*utils.Response, *utils.AppError)
	Remove(id string) (*utils.Response, *utils.AppError)
	Search(page string, limit string, query string) (*utils.Pagination, *utils.AppError)
}

type listingService struct {
	cld         cloudinary.Cloudinary
	listingRepo storage.ListingRepository
	photoRepo   storage.PhotoRepository
	userRepo    storage.UserRepository
	catalogRepo storage.CatalogRepository
}

func NewListingService() ListingService {
	ctx := context.Background()

	db, err := database.GetDatabase(ctx)
	if err != nil {
		log.Msg.Panic("error getting database connection %s", err)
	}

	cld, err := cloudinary.NewCloudinaryService()

	if err != nil {
		log.Msg.DPanicf("error creating cloudinary service: %s", err.Error())
	}

	return &listingService{
		cld:         cld,
		listingRepo: storage.NewListingRepository(db, ctx),
		photoRepo:   storage.NewPhotoRepository(db, ctx),
		userRepo:    storage.NewUserRepository(db, ctx),
		catalogRepo: storage.NewCatalogRepository(db, ctx),
	}
}

func (s *listingService) Search(page string, limit string, query string) (*utils.Pagination, *utils.AppError) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}

	listings, totalItems, totalPage, err := s.listingRepo.SearchByLocation(pageInt, limitInt, query)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	for _, listing := range listings {

		landlord, err := s.userRepo.FindLandlord(listing.LandlordID)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		listing.Landlord = landlord

		photos, err := s.photoRepo.FindAllForListing(listing.ID)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		listing.Photos = photos
	}

	res := utils.NewPaginationResponse(totalItems, totalPage, pageInt, limitInt, listings)

	return res, nil
}

func (s *listingService) FindAll(page string, limit string) (*utils.Pagination, *utils.AppError) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}

	listings, totalItems, totalPage, err := s.listingRepo.FindAll(pageInt, limitInt)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	for _, listing := range listings {

		landlord, err := s.userRepo.FindLandlord(listing.LandlordID)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		listing.Landlord = landlord

		photos, err := s.photoRepo.FindAllForListing(listing.ID)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		listing.Photos = photos
	}

	res := utils.NewPaginationResponse(totalItems, totalPage, pageInt, limitInt, listings)

	return res, nil
}

func (s *listingService) Save(payload *utils.JwtPayload, req *dto.ListingRequest, files []multipart.File) (*utils.Response, error) {

	var catalogs []*domain.Catalog

	for _, catalogId := range req.Catalogs {
		existingCatalog, err := s.catalogRepo.FindById(catalogId)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, utils.NewAppError(404, "Catalog not found!")
			}
			return nil, utils.NewAppError(500, err.Error())
		}

		catalogs = append(catalogs, existingCatalog)
	}

	newListing := &domain.Listing{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Location:    req.Location,
		Guests:      req.Guests,
		Beds:        req.Beds,
		Baths:       req.Baths,
		CleaningFee: req.CleaningFee,
		ServiceFee:  req.ServiceFee,
		Taxes:       req.Taxes,
		LandlordID:  payload.Sub,
		Catalogs:    catalogs,
	}

	listing, err := s.listingRepo.Save(newListing)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	for _, catalog := range catalogs {
		newCatalogListing := &domain.CatalogListing{
			CatalogID: catalog.ID,
			ListingID: listing.ID,
		}

		err = s.catalogRepo.InsertCatalogForListing(newCatalogListing)
		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}
	}

	var photos []*domain.Photo

	for _, file := range files {

		img, err := s.cld.UploadFile(file)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		newPhoto := &domain.Photo{
			ListingID: listing.ID,
			PublicID:  img.PublicID,
			URL:       img.SecureURL,
		}

		photo, err := s.photoRepo.Insert(newPhoto)

		if err != nil {
			return nil, utils.NewAppError(500, err.Error())
		}

		photos = append(photos, photo)
	}

	listing.Photos = photos

	landlord, err := s.userRepo.FindLandlord(listing.LandlordID)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	listing.Landlord = landlord

	return utils.NewResponse(201, listing), nil
}

func (s *listingService) FindDetail(id string) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	listing, err := s.listingRepo.FindOne(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Listing not found!")
		}
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	catalogs, err := s.catalogRepo.FindCatalogsByListingId(listing.ID)

	if err != nil {
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	listing.Catalogs = catalogs

	landlord, err := s.userRepo.FindLandlord(listing.LandlordID)

	if err != nil {
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	listing.Landlord = landlord

	photos, err := s.photoRepo.FindAllForListing(idInt)

	if err != nil {
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	listing.Photos = photos

	return utils.NewResponse(200, listing), nil
}

func (s *listingService) Remove(id string) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	err = s.listingRepo.Remove(idInt)

	if err != nil {
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	return utils.NewResponse(200, "Deleted listing successfully!"), nil
}

func (s *listingService) Update(id string, req *dto.ListingRequest) (*utils.Response, *utils.AppError) {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		return nil, utils.NewAppError(400, "Invalid input")
	}

	existingListing, err := s.listingRepo.FindOne(idInt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Listing not found!")
		}
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	existingListing.Title = req.Title
	existingListing.Description = req.Description
	existingListing.Price = req.Price
	existingListing.Location = req.Location
	existingListing.Guests = req.Guests
	existingListing.Beds = req.Beds
	existingListing.Baths = req.Baths
	existingListing.CleaningFee = req.CleaningFee
	existingListing.ServiceFee = req.ServiceFee
	existingListing.Taxes = req.Taxes

	listing, err := s.listingRepo.Update(idInt, existingListing)

	if err != nil {
		log.Msg.Error(err)
		return nil, utils.NewAppError(500, err.Error())
	}

	return utils.NewResponse(200, listing), nil
}
