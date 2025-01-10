package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/may20xx/booking/config"
	"github.com/may20xx/booking/internal/domain"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JwtPayload struct {
	Sub      int      `json:"sub"`
	Iat      int64    `json:"iat"`
	Exp      int64    `json:"exp"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func ConvertRolesToStrings(roles []domain.Role) []string {
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName)
	}
	return roleNames
}

func GenerateJWT(user *domain.User) (*TokenResponse, error) {
	config := config.GetConfig()
	now := time.Now()

	roleNames := ConvertRolesToStrings(user.Roles)

	accessPayload := JwtPayload{
		Sub:      user.ID,
		Iat:      now.Unix(),
		Exp:      now.Add(30 * 24 * time.Hour).Unix(),
		Username: user.Username,
		Roles:    roleNames,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      accessPayload.Sub,
		"iat":      accessPayload.Iat,
		"exp":      accessPayload.Exp,
		"username": accessPayload.Username,
		"roles":    accessPayload.Roles,
	})

	accessTokenString, err := accessToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return nil, err
	}

	refreshPayload := JwtPayload{
		Sub: user.ID,
		Iat: now.Unix(),
		Exp: now.Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": refreshPayload.Sub,
		"iat": refreshPayload.Iat,
		"exp": refreshPayload.Exp,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWTRefreshSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateJWT(tokenString string) (*JwtPayload, error) {
	config := config.GetConfig()

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid sub claim")
		}

		username, ok := claims["username"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid username claim")
		}

		iat, ok := claims["iat"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid iat claim")
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid exp claim")
		}

		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid roles claim")
		}

		var roles []string
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}

		payload := &JwtPayload{
			Sub:      int(sub),
			Iat:      int64(iat),
			Exp:      int64(exp),
			Username: username,
			Roles:    roles,
		}

		return payload, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
