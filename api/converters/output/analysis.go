package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/domain"
)

func ConvertAnalysisCategory(domainAnalysis *domain.EcosystemAnalysisCategory) *model.EcosystemAnalysisCategory {
	//goland:noinspection GoPreferNilSlice
	messages := []*model.EcosystemAnalysisMessage{}

	for _, domainAnalysisMessage := range domainAnalysis.Messages {
		analysisMessage := ConvertAnalysisMessage(domainAnalysisMessage)
		messages = append(messages, analysisMessage)
	}

	return &model.EcosystemAnalysisCategory{
		ID:          domainAnalysis.ID,
		Name:        domainAnalysis.Name,
		Description: domainAnalysis.Description,
		Status:      model.AnalysisStatus(domainAnalysis.Status),
		Messages:    messages,
	}
}

func ConvertAnalysisMessage(domainAnalysisMessage *domain.EcosystemAnalysisMessage) *model.EcosystemAnalysisMessage {
	return &model.EcosystemAnalysisMessage{
		ID:          domainAnalysisMessage.ID,
		Name:        domainAnalysisMessage.Name,
		Description: domainAnalysisMessage.Description,
		Status:      model.AnalysisStatus(domainAnalysisMessage.Status),
	}
}
