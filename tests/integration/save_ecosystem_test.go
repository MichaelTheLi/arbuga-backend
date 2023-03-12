package integration_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SaveEcosystemResponse struct {
	SaveEcosystem model.EcosystemUpdateResult `json:"saveEcosystem"`
}

func TestSaveEcosystemCreatedEntity(t *testing.T) {
	query := "mutation SaveEcosystem($ecosystem: EcosystemInput!) { saveEcosystem(ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"ecosystem\": { \"name\": \"tEst eCosystem\", \"aquarium\": { \"dimensions\": { \"width\": 10 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &state.Token)
	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems, 1)
	assert.Equal(t, "tEst eCosystem", user.Ecosystems[0].Ecosystem.Name)
	assert.Equal(t, 10, user.Ecosystems[0].Ecosystem.Aquarium.Dimensions.Width)
	assert.Equal(t, 0, user.Ecosystems[0].Ecosystem.Aquarium.Dimensions.Height)
}

func TestSaveEcosystemUpdatesEntity(t *testing.T) {
	query := "mutation SaveEcosystem($id: ID!, $ecosystem: EcosystemInput!) { saveEcosystem(id: $id, ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"id\": \"testId1\",\"ecosystem\": { \"name\": \"tEst eCosystem Updated\", \"aquarium\": { \"dimensions\": { \"width\": 11, \"height\": 12, \"length\": 13 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	oldUser, _ := state.State.UserGateway.GetUserByLogin("testLogin")
	addTestEcosystem(oldUser)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &state.Token)
	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems, 1)
	assert.Equal(t, "tEst eCosystem Updated", user.Ecosystems[0].Ecosystem.Name)
	assert.Equal(t, 11, user.Ecosystems[0].Ecosystem.Aquarium.Dimensions.Width)
	assert.Equal(t, 12, user.Ecosystems[0].Ecosystem.Aquarium.Dimensions.Height)
	assert.Equal(t, 13, user.Ecosystems[0].Ecosystem.Aquarium.Dimensions.Length)
}

func TestCantSaveOrUpdateEcosystemIfNotAuthenticated(t *testing.T) {
	query := "mutation SaveEcosystem($id: ID!, $ecosystem: EcosystemInput!) { saveEcosystem(id: $id, ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"id\": \"testId1\",\"ecosystem\": { \"name\": \"tEst eCosystem\", \"aquarium\": { \"dimensions\": { \"width\": 10 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	oldUser, _ := state.State.UserGateway.GetUserByLogin("testLogin")
	addTestEcosystem(oldUser)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, nil)
	user, userErr := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems, 1)
	assert.Equal(t, "Old Name", user.Ecosystems[0].Ecosystem.Name)

	err := data.Errors[0]

	assert.Equal(t, "saveEcosystem", err.Path.String())
	assert.Equal(t, "not authenticated", err.Message)
}

func TestUpdateInvalidIdWillError(t *testing.T) {
	query := "mutation SaveEcosystem($id: ID!, $ecosystem: EcosystemInput!) { saveEcosystem(id: $id, ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"id\": \"testId2NotExist\",\"ecosystem\": { \"name\": \"tEst eCosystem\", \"aquarium\": { \"dimensions\": { \"width\": 10 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	oldUser, _ := state.State.UserGateway.GetUserByLogin("testLogin")
	addTestEcosystem(oldUser)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &state.Token)
	user, userErr := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems, 1)
	assert.Equal(t, "Old Name", user.Ecosystems[0].Ecosystem.Name)

	var saveData SaveEcosystemResponse
	err := json.Unmarshal(data.Data, &saveData)
	assert.Nil(t, err)
	assert.False(t, saveData.SaveEcosystem.Success)
	assert.Nil(t, saveData.SaveEcosystem.Ecosystem)
	assert.Equal(t, "Ecosystem not found", *saveData.SaveEcosystem.Error)
}

func addTestEcosystem(oldUser *app.User) {
	oldUser.Ecosystems = append(oldUser.Ecosystems, &app.Ecosystem{
		ID: "testId1",
		Ecosystem: &domain.Ecosystem{
			Name: "Old Name",
			Aquarium: &domain.AquariumGlass{
				Dimensions: &domain.Dimensions{
					Width:  10,
					Height: 10,
					Length: 10,
				},
				GlassThickness:     1,
				SubstrateThickness: nil,
				DecorationsVolume:  nil,
			},
			Analysis: nil,
		},
	})
}
