package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/domain"
)

func ConvertEcosystem(domainEcosystem *domain.Ecosystem) *model.Ecosystem {
	var analysis []*model.EcosystemAnalysisCategory

	for _, domainAnalysisCategory := range domainEcosystem.Analysis {
		analysisCategory := ConvertAnalysisCategory(domainAnalysisCategory)
		analysis = append(analysis, analysisCategory)
	}

	var aquarium *model.AquariumGlass
	if domainEcosystem.Aquarium != nil {
		aquarium = ConvertAquarium(domainEcosystem.Aquarium)
	}

	return &model.Ecosystem{
		ID:       domainEcosystem.ID,
		Name:     domainEcosystem.Name,
		Aquarium: aquarium,
		Analysis: analysis,
	}
}
