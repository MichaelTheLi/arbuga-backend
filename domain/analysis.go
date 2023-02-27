package domain

type EcosystemAnalysisCategory struct {
	ID          string
	Name        string
	Description string
	Status      AnalysisStatus
	Messages    []*EcosystemAnalysisMessage
}

type EcosystemAnalysisMessage struct {
	ID          string
	Name        string
	Description string
	Status      AnalysisStatus
}

type AnalysisStatus string
