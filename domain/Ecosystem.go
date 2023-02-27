package domain

type Ecosystem struct {
	ID       string
	Name     string
	Aquarium *AquariumGlass
	Analysis []*EcosystemAnalysisCategory
}
