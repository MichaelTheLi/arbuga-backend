package adapters

import (
	"arbuga/backend/domain"
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type JwtTokenService struct {
	Secret string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId string
	Name   string
}

func (service JwtTokenService) GenerateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), UserClaims{
		UserId:           user.ID,
		Name:             user.Name,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	tokenString, err := token.SignedString([]byte(service.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service JwtTokenService) GetUserIdFromToken(tokenValue string) (string, error) {
	userClaims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.Secret), nil
	})

	// Checking token validity
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return userClaims.UserId, nil
}
