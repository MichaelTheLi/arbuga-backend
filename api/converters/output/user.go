package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
)

func ConvertUser(domainUser *app.User) *model.User {
	var ecosystems []*model.Ecosystem

	for _, ecosystem := range domainUser.Ecosystems {
		newEcosystem := ConvertEcosystem(ecosystem)
		ecosystems = append(ecosystems, newEcosystem)
	}

	return &model.User{
		ID:         domainUser.ID,
		Login:      domainUser.Login,
		Name:       domainUser.Owner.Name,
		Ecosystems: ecosystems,
	}
}
