package app

import (
	"arbuga/backend/domain"
	"errors"
)

type SignInService struct {
	Gateway      UserGateway
	AuthService  AuthService
	TokenService TokenService
}

type LoginResult struct {
	Token string
	User  *domain.User
}

func (service *SignInService) Login(login string, password string) (*LoginResult, error) {
	user, _ := service.Gateway.GetUserByLogin(login)

	if user == nil {
		return nil, errors.New("unknown user")
	}

	valid, err := service.AuthService.IsHashValid(*user.PasswordHash, password)
	if valid == false || err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := service.TokenService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("token cannot be generated")
	}

	return &LoginResult{User: user, Token: token}, nil
}
