package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/may20xx/booking/internal/api/dto"
	"github.com/may20xx/booking/internal/database"
	"github.com/may20xx/booking/internal/domain"
	"github.com/may20xx/booking/internal/storage"
	"github.com/may20xx/booking/internal/utils"
	"github.com/may20xx/booking/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler interface {
	RegisterHandler(request *dto.RegisterRequest) (*domain.User, *utils.AppError)
	LoginHandler(request *dto.LoginRequest) (*utils.TokenResponse, *utils.AppError)
}

type AuthService struct {
	userRepo  storage.UserStorage
	tokenRepo storage.TokenStorage
	roleRepo  storage.RoleStorage
}

func NewAuthService() *AuthService {
	ctx := context.Background()

	db, err := database.GetDatabase(ctx)
	if err != nil {
		log.Msg.Panic("error getting database connection %s", err)
	}

	return &AuthService{
		userRepo:  storage.NewUserRepository(db, ctx),
		tokenRepo: storage.NewTokenRepository(db, ctx),
		roleRepo:  storage.NewRoleRepository(db, ctx),
	}
}

func (s *AuthService) RegisterHandler(request *dto.RegisterRequest) (*domain.User, *utils.AppError) {
	existingUser, _ := s.userRepo.FindOneByEmail(request.Email)

	if existingUser != nil {
		return nil, utils.NewAppError(400, "User already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	user := &domain.User{
		Email:        request.Email,
		Avatar:       nil,
		FirstName:    request.FirstName,
		Username:     request.Username,
		HashPassword: string(hash),
		Surname:      request.Surname,
	}

	userRole, err := s.roleRepo.FindRoleByName(dto.User)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "User role not found")
		}
		return nil, utils.NewAppError(500, "Internal server error")
	}

	user.Roles = append(user.Roles, *userRole)

	result, err := s.userRepo.Insert(user)

	if err != nil {
		return nil, utils.NewAppError(500, err.Error())
	}

	return result, nil
}

func (s *AuthService) LoginHandler(request *dto.LoginRequest) (*utils.TokenResponse, *utils.AppError) {
	user, err := s.userRepo.FindOneByUsername(request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewAppError(404, "Username or password is incorrect")
		}
		return nil, utils.NewAppError(500, "Internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(request.Password))
	if err != nil {
		return nil, utils.NewAppError(401, "Username or password is incorrect")
	}

	roles, _ := s.roleRepo.FindRolesByUser(user.ID)

	user.Roles = roles

	var refreshToken string
	var needNewRefreshToken bool = true

	existingToken, err := s.tokenRepo.FindOneByToken(dto.RefreshToken, user.ID)
	if err != nil {
		return nil, utils.NewAppError(500, "Error checking existing token")
	}

	if existingToken != nil && existingToken.ExpiredAt != nil && existingToken.ExpiredAt.After(time.Now()) {
		refreshToken = existingToken.Token
		needNewRefreshToken = false
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, utils.NewAppError(400, "Error generating JWT")
	}

	if needNewRefreshToken {
		refreshToken = token.RefreshToken
		expirationTime := time.Now().Add(30 * 24 * time.Hour)

		tokenEntity := &domain.Token{
			UserID:    user.ID,
			Name:      dto.RefreshToken,
			Token:     refreshToken,
			ExpiredAt: &expirationTime,
		}

		if existingToken != nil {
			tokenEntity.ID = existingToken.ID
			_, err = s.tokenRepo.Update(tokenEntity)
		} else {
			_, err = s.tokenRepo.Insert(tokenEntity)
		}

		if err != nil {
			return nil, utils.NewAppError(500, "Error saving refresh token")
		}
	}

	tokenResponse := &utils.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: refreshToken,
	}

	return tokenResponse, nil
}
