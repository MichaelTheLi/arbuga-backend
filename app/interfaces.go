package app

import (
	"arbuga/backend/domain"
)

type UserGateway interface {
	GetUserByLogin(login string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	CreateUser(login string, hash string, name string) (*domain.User, error)
}

type AuthService interface {
	HashFromPassword(password string) (string, error)
	IsHashValid(expectedHash string, password string) (bool, error)
}

type TokenService interface {
	GenerateToken(user *domain.User) (string, error)
	GetUserIdFromToken(token string) (string, error)
}
