package integration_test

import (
	"arbuga/backend/api/graph/model"
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

	assert.Len(t, user.Owner.Ecosystems, 1)
	assert.Equal(t, "tEst eCosystem", user.Owner.Ecosystems[0].Name)
	assert.Equal(t, 10, user.Owner.Ecosystems[0].Aquarium.Dimensions.Width)
	assert.Equal(t, 0, user.Owner.Ecosystems[0].Aquarium.Dimensions.Height)
}

func TestSaveEcosystemUpdatesEntity(t *testing.T) {
	query := "mutation SaveEcosystem($id: ID!, $ecosystem: EcosystemInput!) { saveEcosystem(id: $id, ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"id\": \"testId1\",\"ecosystem\": { \"name\": \"tEst eCosystem Updated\", \"aquarium\": { \"dimensions\": { \"width\": 11, \"height\": 12, \"length\": 13 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	oldUser, _ := state.State.UserGateway.GetUserByLogin("testLogin")
	oldUser.Owner.Ecosystems = append(oldUser.Owner.Ecosystems, &domain.Ecosystem{
		ID:   "testId1",
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
	})

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &state.Token)
	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Owner.Ecosystems, 1)
	assert.Equal(t, "tEst eCosystem Updated", user.Owner.Ecosystems[0].Name)
	assert.Equal(t, 11, user.Owner.Ecosystems[0].Aquarium.Dimensions.Width)
	assert.Equal(t, 12, user.Owner.Ecosystems[0].Aquarium.Dimensions.Height)
	assert.Equal(t, 13, user.Owner.Ecosystems[0].Aquarium.Dimensions.Length)
}

func TestCantSaveOrUpdateEcosystemIfNotAuthenticated(t *testing.T) {
	query := "mutation SaveEcosystem($id: ID!, $ecosystem: EcosystemInput!) { saveEcosystem(id: $id, ecosystem: $ecosystem) { success error ecosystem { id name aquarium {dimensions{width height} }}}}"
	//goland:noinspection SpellCheckingInspection
	variables := "{ \"id\": \"testId1\",\"ecosystem\": { \"name\": \"tEst eCosystem\", \"aquarium\": { \"dimensions\": { \"width\": 10 } } } }"

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")
	oldUser, _ := state.State.UserGateway.GetUserByLogin("testLogin")
	oldUser.Owner.Ecosystems = append(oldUser.Owner.Ecosystems, &domain.Ecosystem{
		ID:   "testId1",
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
	})

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, nil)
	user, userErr := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	assert.Len(t, user.Owner.Ecosystems, 1)
	assert.Equal(t, "Old Name", user.Owner.Ecosystems[0].Name)

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
	oldUser.Owner.Ecosystems = append(oldUser.Owner.Ecosystems, &domain.Ecosystem{
		ID:   "testId1",
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
	})

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "SaveEcosystem", &data, &state.Token)
	user, userErr := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	assert.Len(t, user.Owner.Ecosystems, 1)
	assert.Equal(t, "Old Name", user.Owner.Ecosystems[0].Name)

	var saveData SaveEcosystemResponse
	err := json.Unmarshal(data.Data, &saveData)
	assert.Nil(t, err)
	assert.False(t, saveData.SaveEcosystem.Success)
	assert.Nil(t, saveData.SaveEcosystem.Ecosystem)
	assert.Equal(t, "Ecosystem not found", *saveData.SaveEcosystem.Error)
}
