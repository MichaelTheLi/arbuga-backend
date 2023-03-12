package integration_test

import (
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyUserDoesntHaveFish(t *testing.T) {
	_, _, meData := executeGetMe(t, false)
	assert.Empty(t, meData.Me.Ecosystems[0].Fish)
}

func TestMeReturnsFishList(t *testing.T) {
	_, _, meData := executeGetMe(t, true)
	assert.Len(t, meData.Me.Ecosystems[0].Fish, 2)
}

func TestMeReturnsCorrectFishList(t *testing.T) {
	_, _, meData := executeGetMe(t, true)
	assert.Equal(t, meData.Me.Ecosystems[0].Fish[0].ID, "test1")
	assert.Equal(t, meData.Me.Ecosystems[0].Fish[0].Name, "Test fish 1")
	assert.Equal(t, meData.Me.Ecosystems[0].Fish[1].ID, "test2")
	assert.Equal(t, meData.Me.Ecosystems[0].Fish[1].Name, "Test fish 2")
}

func executeGetMe(t *testing.T, withFish bool) (utils.TestServerState, graphql.Response, MeResponse) {
	query := "query Me {me {id ecosystems {fish {id name}}}}"
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

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)
	return state, data, meData
}
