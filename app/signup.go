package app

import (
	"arbuga/backend/domain"
	"errors"
)

type SignUpService struct {
	Gateway     UserGateway
	AuthService AuthService
}

type SignUpResult struct {
	User *domain.User
}

func (service *SignUpService) SignUp(login string, password string, name string) (*SignUpResult, error) {
	user, _ := service.Gateway.GetUserByLogin(login)

	if user != nil {
		return nil, errors.New("already exists")
	}

	hashedPassString, err := service.AuthService.HashFromPassword(password)

	if err != nil {
		return nil, err
	}

	newUser, errCreate := service.Gateway.CreateUser(login, hashedPassString, name)

	if errCreate != nil {
		return nil, errCreate
	}

	return &SignUpResult{User: newUser}, nil
}
