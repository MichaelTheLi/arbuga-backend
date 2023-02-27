package adapters

import (
	"arbuga/backend/domain"
	"crypto/rand"
	"fmt"
	"math/big"
)

type LocalUsersGateway struct {
	Users map[string]*domain.User
}

func (state LocalUsersGateway) GetUserByLogin(login string) (*domain.User, error) {
	var user *domain.User
	for _, v := range state.Users {
		if v.Login != nil && *v.Login == login {
			user = v
			break
		}
	}

	return user, nil
}

func (state LocalUsersGateway) GetUserByID(id string) (*domain.User, error) {
	user, _ := state.Users[id]
	return user, nil
}

func (state LocalUsersGateway) CreateUser(login string, hash string, name string) (*domain.User, error) {
	randValue, _ := rand.Int(rand.Reader, big.NewInt(100))
	user := &domain.User{
		ID:           fmt.Sprintf("T%d", randValue),
		Name:         name,
		Login:        &login,
		PasswordHash: &hash,
		Ecosystems:   []*domain.Ecosystem{},
	}
	state.Users[user.ID] = user

	return user, nil
}
