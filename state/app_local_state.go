package state

import (
	"arbuga/backend/graph/model"
)

type AppLocalState struct {
	Users map[string]*model.User
}

func (state AppLocalState) GetUserByID(id string) (*model.User, error) {
	user, _ := state.Users[id]
	return user, nil
}

func (state AppLocalState) GetUserByLogin(login string) (*model.User, error) {
	var user *model.User
	for _, v := range state.Users {
		if v.Login != nil && *v.Login == login {
			user = v
			break
		}
	}

	return user, nil
}
