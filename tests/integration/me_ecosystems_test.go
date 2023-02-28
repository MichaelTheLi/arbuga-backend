package integration_test

import (
	"arbuga/backend/api/converters/output"
	"arbuga/backend/domain"
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcosystemsValidCount(t *testing.T) {
	query := "query Me {me {id ecosystems { id }}}"
	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	state.UsersGateway.Users["testId"].Owner.Ecosystems = []*domain.Ecosystem{
		{
			ID: "ecosystem1",
		},
		{
			ID: "ecosystem2",
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Len(t, meData.Me.Ecosystems, len(state.UsersGateway.Users["testId"].Owner.Ecosystems))
}

func TestEcosystemsBasicDataReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id name }}}"
	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	state.UsersGateway.Users["testId"].Owner.Ecosystems = []*domain.Ecosystem{
		{
			ID:   "ecosystem1",
			Name: "Test Ecosystem Name",
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	convertedUser := output.ConvertUser(state.UsersGateway.Users["testId"])
	assert.Equal(t, convertedUser.Ecosystems[0], meData.Me.Ecosystems[0])
}

func TestEcosystemsAquariumDimensionsReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id aquarium {glassThickness substrateThickness decorationsVolume dimensions{width height length}} }}}"
	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	thickness := 15
	volume := 18
	state.UsersGateway.Users["testId"].Owner.Ecosystems = []*domain.Ecosystem{
		{
			ID: "ecosystem1",
			Aquarium: &domain.AquariumGlass{
				Dimensions: &domain.Dimensions{
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

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	convertedUser := output.ConvertUser(state.UsersGateway.Users["testId"])
	assert.Len(t, meData.Me.Ecosystems, 1)
	assert.NotNil(t, meData.Me.Ecosystems[0].Aquarium)
	assert.Equal(t, convertedUser.Ecosystems[0].Aquarium, meData.Me.Ecosystems[0].Aquarium)
}

func TestEcosystemsAnalysisReceived(t *testing.T) {
	query := "query Me {me {id ecosystems { id analysis { id name description status messages { id name description status } } }}}"
	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	state.UsersGateway.Users["testId"].Owner.Ecosystems = []*domain.Ecosystem{
		{
			ID: "ecosystem1",
			Analysis: []*domain.EcosystemAnalysisCategory{
				{
					ID:          "test1Category",
					Name:        "CategoryOne",
					Description: "TestDescr1",
					Status:      "ok",
					Messages: []*domain.EcosystemAnalysisMessage{
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
					Messages:    []*domain.EcosystemAnalysisMessage{},
				},
			},
		},
	}

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	convertedUser := output.ConvertUser(state.UsersGateway.Users["testId"])
	assert.Equal(t, convertedUser.Ecosystems[0].Analysis, meData.Me.Ecosystems[0].Analysis)
}
