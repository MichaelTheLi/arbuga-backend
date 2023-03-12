package integration_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AddResponse struct {
	Add model.AddFishResult `json:"addFishToEcosystem"`
}

func TestCanAddFishIntoEmptyEcosystem(t *testing.T) {
	state, data, addData := executeAddFish(t, "ecosystem1", "test3", false)
	assert.NotEmpty(t, data)
	assert.Len(t, addData.Add.Ecosystem.Fish, 1)
	assert.Equal(t, addData.Add.Ecosystem.Fish[0].ID, "test3")

	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems[0].Fish, 1)
	assert.Equal(t, user.Ecosystems[0].Fish[0].Id, "test3")
}

func TestCanAddFishIntoExistingEcosystem(t *testing.T) {
	state, _, addData := executeAddFish(t, "ecosystem1", "test3", true)
	assert.Len(t, addData.Add.Ecosystem.Fish, 3)
	assert.Equal(t, addData.Add.Ecosystem.Fish[0].ID, "test1")
	assert.Equal(t, addData.Add.Ecosystem.Fish[1].ID, "test2")
	assert.Equal(t, addData.Add.Ecosystem.Fish[2].ID, "test3")
	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems[0].Fish, 3)
	assert.Equal(t, user.Ecosystems[0].Fish[0].Id, "test1")
	assert.Equal(t, user.Ecosystems[0].Fish[1].Id, "test2")
	assert.Equal(t, user.Ecosystems[0].Fish[2].Id, "test3")
}

func TestCannotAddFishIntoInvalidEcosystem(t *testing.T) {
	state, data, _ := executeAddFish(t, "invalidId", "test3", true)
	assert.Len(t, data.Errors, 1)
	assert.Equal(t, data.Errors[0].Path.String(), "addFishToEcosystem")
	assert.Equal(t, data.Errors[0].Message, "invalid ecosystem")

	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems[0].Fish, 2)
}

func TestCannotAddInvalidFish(t *testing.T) {
	state, data, _ := executeAddFish(t, "ecosystem1", "invalidId", true)
	assert.Len(t, data.Errors, 1)
	assert.Equal(t, data.Errors[0].Path.String(), "addFishToEcosystem")
	assert.Equal(t, data.Errors[0].Message, "invalid fish")

	user, err := state.State.UserGateway.GetUserByLogin("testLogin")

	assert.Nil(t, err)
	assert.NotNil(t, user)

	assert.Len(t, user.Ecosystems[0].Fish, 2)
}

func executeAddFish(t *testing.T, ecosystemId string, fishId string, withFish bool) (utils.TestServerState, graphql.Response, AddResponse) {
	query := "mutation Add($ecosystemId: ID!, $fishId: ID!) {addFishToEcosystem(ecosystemId: $ecosystemId, fishId: $fishId) {ecosystem {id fish {id name}}}}"
	variables := fmt.Sprintf("{ \"ecosystemId\":  \"%s\",  \"fishId\":  \"%s\"}", ecosystemId, fishId)

	var data graphql.Response
	stateWithUser := utils.BuildStateWithUser("testLogin", "testPass")
	state := utils.BuildStateWithFish(&stateWithUser)

	state.UsersGateway.Users["testId"].Ecosystems = []*app.Ecosystem{
		{
			ID: "ecosystem1",
			Ecosystem: &domain.Ecosystem{
				Name: "Ecosystem 1",
			},
		},
	}

	if withFish {
		state.UsersGateway.Users["testId"].Ecosystems[0].Fish = []*app.Fish{
			{
				Id: "test1",
				Fish: &domain.Fish{
					Name:        "Test fish 1",
					Description: "Test desc 1",
				},
			},
			{
				Id: "test2",
				Fish: &domain.Fish{
					Name:        "Test fish 2",
					Description: "Test desc 2",
				},
			},
		}
	}

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "Add", &data, &state.Token)

	var addData AddResponse
	err := json.Unmarshal(data.Data, &addData)
	assert.Nil(t, err)
	return state, data, addData
}
