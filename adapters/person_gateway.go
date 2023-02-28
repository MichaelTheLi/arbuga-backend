package adapters

import (
	"arbuga/backend/app"
)

type LocalUserGateway struct {
	Users map[string]*app.User
}

func (state LocalUserGateway) GetUserByLogin(login string) (*app.User, error) {
	var user *app.User
	for _, v := range state.Users {
		if v.Login != nil && *v.Login == login {
			user = v
			break
		}
	}

	return user, nil
}

func (state LocalUserGateway) GetUserByID(id string) (*app.User, error) {
	user, _ := state.Users[id]
	return user, nil
}

func (state LocalUserGateway) SaveUser(user *app.User) (*app.User, error) {
	state.Users[user.ID] = user

	return user, nil
}
