package app

import (
	"arbuga/backend/domain"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type SignUpService struct {
	Gateway     UserGateway
	AuthService AuthService
}

type SignUpResult struct {
	User *User
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

	randValue, _ := rand.Int(rand.Reader, big.NewInt(100))
	newOwner := domain.NewOwner(name)
	newUser := &User{
		ID:           fmt.Sprintf("T%d", randValue),
		Owner:        newOwner,
		Login:        &login,
		PasswordHash: &hashedPassString,
	}
	SavedUser, errCreate := service.Gateway.SaveUser(newUser)

	if errCreate != nil {
		return nil, errCreate
	}

	return &SignUpResult{User: SavedUser}, nil
}
