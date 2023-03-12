package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
)

func ConvertEcosystem(domainEcosystem *app.Ecosystem) *model.Ecosystem {
	var analysis []*model.EcosystemAnalysisCategory

	if domainEcosystem.Ecosystem.Aquarium != nil {
		for _, domainAnalysisCategory := range domainEcosystem.Ecosystem.Analysis {
			analysisCategory := ConvertAnalysisCategory(domainAnalysisCategory)
			analysis = append(analysis, analysisCategory)
		}
	}

	var aquarium *model.AquariumGlass
	if domainEcosystem.Ecosystem.Aquarium != nil {
		aquarium = ConvertAquarium(domainEcosystem.Ecosystem.Aquarium)
	}

	return &model.Ecosystem{
		ID:       domainEcosystem.ID,
		Name:     domainEcosystem.Ecosystem.Name,
		Aquarium: aquarium,
		Analysis: analysis,
		Fish:     ConvertFishList(domainEcosystem.Fish),
	}
}
