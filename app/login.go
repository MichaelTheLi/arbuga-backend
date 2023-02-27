package app

import (
	"arbuga/backend/domain"
	"errors"
)

type LoginService struct {
	gateway     UserGateway
	authService AuthService
}

type LoginResult struct {
	Token string
	User  *domain.User
}

func (service *LoginService) Login(login string, password string) (*LoginResult, error) {
	user, _ := service.gateway.GetUserByLogin(login)

	if user == nil {
		return nil, errors.New("unknown user")
	}

	valid, err := service.authService.IsHashValid(*user.PasswordHash, password)
	if valid == false || err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := service.authService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("token cannot be generated")
	}

	return &LoginResult{User: user, Token: token}, nil
}
