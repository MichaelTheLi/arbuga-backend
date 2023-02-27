package app

import (
	"arbuga/backend/domain"
)

type UserGateway interface {
	GetUserByLogin(login string) (*domain.User, error)
	CreateUser(login string, hash string, token string) (*domain.User, error)
}

type AuthService interface {
	HashFromPassword(password string) (string, error)
	IsHashValid(expectedHash string, password string) (bool, error)
	GenerateToken(user *domain.User) (string, error)
}
