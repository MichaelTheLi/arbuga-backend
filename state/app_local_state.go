package state

import (
	"arbuga/backend/graph/model"
	"errors"
)

type AppLocalState struct {
	Users map[string]*model.User
}

func (state AppLocalState) GetUserByID(id string) (*model.User, error) {
	user, _ := state.Users[id]
	return user, nil
}

func (state AppLocalState) GetUserByLoginAndPassword(login, password string) (*model.User, error) {
	var user *model.User
	for _, v := range state.Users {
		if v.Login != nil && v.Password != nil && *v.Login == login {
			if *v.Password == password { // TODO Yeah, you know..
				user = v
			} else {
				return nil, errors.New("invalid credentials")
			}
			break
		}
	}

	return user, nil
}
