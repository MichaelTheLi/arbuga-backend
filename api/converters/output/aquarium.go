package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/domain"
)

func ConvertAquarium(domainAquarium *domain.AquariumGlass) *model.AquariumGlass {
	dimensions := model.Dimensions(*domainAquarium.Dimensions)

	return &model.AquariumGlass{
		Dimensions:         &dimensions,
		GlassThickness:     domainAquarium.GlassThickness,
		SubstrateThickness: domainAquarium.SubstrateThickness,
		DecorationsVolume:  domainAquarium.DecorationsVolume,
	}
}
