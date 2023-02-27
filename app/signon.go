package app

import (
	"arbuga/backend/domain"
	"errors"
)

type SignOnService struct {
	gateway     UserGateway
	authService AuthService
}

type SignOnResult struct {
	User *domain.User
}

func (service *SignOnService) SignOn(login string, password string, name string) (*SignOnResult, error) {
	user, _ := service.gateway.GetUserByLogin(login)

	if user != nil {
		return nil, errors.New("already exists")
	}

	hashedPassString, err := service.authService.HashFromPassword(password)

	if err != nil {
		return nil, err
	}

	newUser, errCreate := service.gateway.CreateUser(login, hashedPassString, name)

	if errCreate != nil {
		return nil, errCreate
	}

	return &SignOnResult{User: newUser}, nil
}
