package domain

type Ecosystem struct {
	ID       string
	Name     string
	Aquarium *AquariumGlass
	Analysis []*EcosystemAnalysisCategory
}

// CalculateAnalysis TODO Should it be domain?
func (ecosystem *Ecosystem) CalculateAnalysis() []*EcosystemAnalysisCategory {
	return []*EcosystemAnalysisCategory{
		{
			ID:          "1",
			Name:        "Filtration",
			Description: "How good filtration is",
			Status:      "ok",
			Messages:    nil,
		},
		{
			ID:          "2",
			Name:        "Temperature",
			Description: "How good filtration is",
			Status:      "moderate",
			Messages: []*EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Low temperature: fish",
					Description: "Temperature too low for some of the fish",
					Status:      "bad",
				},
				{
					ID:          "2",
					Name:        "Not optimal temperature: fish",
					Description: "Temperature on the low edge for some of the fish",
					Status:      "moderate",
				},
			},
		},
		{
			ID:          "3",
			Name:        "Fish compatibility",
			Description: "Are your fish compatible with each other",
			Status:      "bad",
			Messages: []*EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Barbus",
					Description: "These guys are likely to behave aggressive towards most of your fish",
					Status:      "bad",
				},
				{
					ID:          "2",
					Name:        "Bolivian Ram",
					Description: "Bolivian Ram are too shy for the neighbourhood with barbuses, you might notice some issues with the food",
					Status:      "moderate",
				},
			},
		},
		{
			ID:          "3",
			Name:        "Plants compatibility",
			Description: "Are your plants compatible with the water in your tank",
			Status:      "moderate",
			Messages: []*EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Cryptocoryne might melt",
					Description: "Cryptocoryne might melt, water too soft",
					Status:      "moderate",
				},
			},
		},
	}
}
