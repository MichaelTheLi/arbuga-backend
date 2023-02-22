package main_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/tests/utils"
	json "encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcosystemsValidCount(t *testing.T) {
	query := "query Me {me {id ecosystems { id }}}"
	var data graphql.Response
	state, token := utils.BuildStateWithUser("testLogin", "testPass")
	state.Users["testId"].Ecosystems = []*model.Ecosystem{
		{
			ID: "ecosystem1",
		},
		{
			ID: "ecosystem2",
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Len(t, meData.Me.Ecosystems, len(state.Users["testId"].Ecosystems))
}

func TestEcosystemsBasicDataReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id name }}}"
	var data graphql.Response
	state, token := utils.BuildStateWithUser("testLogin", "testPass")
	state.Users["testId"].Ecosystems = []*model.Ecosystem{
		{
			ID:   "ecosystem1",
			Name: "Test Ecosystem Name",
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Equal(t, meData.Me.Ecosystems[0], state.Users["testId"].Ecosystems[0])
}

func TestEcosystemsAquariumDimensionsReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id aquarium {glassThickness substrateThickness decorationsVolume dimensions{width height length}} }}}"
	var data graphql.Response
	state, token := utils.BuildStateWithUser("testLogin", "testPass")
	thickness := 15
	volume := 18
	state.Users["testId"].Ecosystems = []*model.Ecosystem{
		{
			ID: "ecosystem1",
			Aquarium: &model.AquariumGlass{
				Dimensions: &model.Dimensions{
					Width:  11,
					Height: 16,
					Length: 21,
				},
				GlassThickness:     12,
				SubstrateThickness: &thickness,
				DecorationsVolume:  &volume,
			},
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Equal(t, meData.Me.Ecosystems[0].Aquarium, state.Users["testId"].Ecosystems[0].Aquarium)
}

func TestEcosystemsAnalysisReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id analysis { id name description status messages { id name description status } } }}}"
	var data graphql.Response
	state, token := utils.BuildStateWithUser("testLogin", "testPass")
	state.Users["testId"].Ecosystems = []*model.Ecosystem{
		{
			ID: "ecosystem1",
			Analysis: []*model.EcosystemAnalysisCategory{
				{
					ID:          "test1Category",
					Name:        "CategoryOne",
					Description: "TestDescr1",
					Status:      "ok",
					Messages: []*model.EcosystemAnalysisMessage{
						{
							ID:          "test1Message",
							Name:        "Test1Msg",
							Description: "TestDescr1Msg",
							Status:      "bad",
						},
						{
							ID:          "test2Message",
							Name:        "Test2Msg",
							Description: "TestDescr2Msg",
							Status:      "moderate",
						},
					},
				},
				{
					ID:          "test2Category",
					Name:        "CategoryTwo",
					Description: "TestDescr2",
					Status:      "moderate",
					Messages:    []*model.EcosystemAnalysisMessage{},
				},
			},
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Equal(t, meData.Me.Ecosystems[0].Analysis, state.Users["testId"].Ecosystems[0].Analysis)
}
