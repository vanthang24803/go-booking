package handler

import (
	"context"
	"database/sql"
	"mime/multipart"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/domain"
	"github.com/may20xx/booking/internal/storage"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/cloudinary"
	"github.com/may20xx/booking/pkg/log"
)

type MeService interface {
	Logout(payload *utils.JwtPayload) *utils.AppError
	GetProfile(payload *utils.JwtPayload) (*domain.User, *utils.AppError)
	UploadAvatar(payload *utils.JwtPayload, file multipart.File) (*domain.User, *utils.AppError)
	UpdateProfile(payload *utils.JwtPayload, req *dto.UpdateProfileRequest) (*domain.User, *utils.AppError)
}

type meService struct {
	userRepo   storage.UserRepository
	roleRepo   storage.RoleStorage
	tokenRepo  storage.TokenStorage
	cloudinary cloudinary.Cloudinary
}

func NewMeService() *meService {
	ctx := context.Background()

	db, err := database.GetDatabase(ctx)

	if err != nil {
		log.Msg.Panic("error getting database connection %s", err)
	}

	upload, err := cloudinary.NewCloudinaryService()

	if err != nil {
		log.Msg.DPanicf("error creating cloudinary service: %s", err.Error())
	}

	return &meService{
		userRepo:   storage.NewUserRepository(db, ctx),
		roleRepo:   storage.NewRoleRepository(db, ctx),
		tokenRepo:  storage.NewTokenRepository(db, ctx),
		cloudinary: upload,
	}
}

func (s *meService) Logout(payload *utils.JwtPayload) *utils.AppError {

	token, _ := s.tokenRepo.FindOneByToken(dto.RefreshToken, payload.Sub)

	if token != nil {
		err := s.tokenRepo.Remove(token.ID)

		if err != nil {
			return utils.NewAppError(500, err.Error())
		}
	}

	return nil
}

func (s *meService) GetProfile(payload *utils.JwtPayload) (*domain.User, *utils.AppError) {
	user, err := s.userRepo.FindOneById(payload.Sub)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, "Internal server error")
	}

	roles, _ := s.roleRepo.FindRolesByUser(user.ID)

	user.Roles = roles

	return user, nil
}

func (s *meService) UploadAvatar(payload *utils.JwtPayload, file multipart.File) (*domain.User, *utils.AppError) {
	user, err := s.userRepo.FindOneById(payload.Sub)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, "Internal server error")
	}

	avt, err := s.cloudinary.UploadFile(file)

	if err != nil {
		return nil, utils.NewAppError(500, "Upload avatar failed")
	}

	user.Avatar = &avt.SecureURL

	user, err = s.userRepo.Update(user)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	return user, nil
}

func (s *meService) UpdateProfile(payload *utils.JwtPayload, req *dto.UpdateProfileRequest) (*domain.User, *utils.AppError) {
	user, err := s.userRepo.FindOneById(payload.Sub)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, "Internal server error")
	}

	user.FirstName = req.FirstName
	user.Surname = req.Surname
	user.Email = req.Email

	res, err := s.userRepo.Update(user)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	return res, nil

}
