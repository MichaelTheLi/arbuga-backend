package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/domain"
)

func ConvertUser(domainUser *domain.User) *model.User {
	var ecosystems []*model.Ecosystem

	for _, ecosystem := range domainUser.Ecosystems {
		newEcosystem := ConvertEcosystem(ecosystem)
		ecosystems = append(ecosystems, newEcosystem)
	}

	return &model.User{
		ID:         domainUser.ID,
		Login:      domainUser.Login,
		Password:   domainUser.PasswordHash,
		Name:       domainUser.Name,
		Ecosystems: ecosystems,
		Ecosystem:  nil,
	}
}
