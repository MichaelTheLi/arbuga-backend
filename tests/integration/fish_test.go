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

type FishResponse struct {
	Fish []*model.Fish `json:"fish"`
}

func TestFishRequestWillReceiveData(t *testing.T) {
	state, _, fishListData, err := executeGetFish(t)
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, len(state.FishGateway.Fish))
}

func executeGetFish(t *testing.T) (utils.TestServerState, graphql.Response, FishResponse, error) {
	query := "query Fish {fish {id name description}}"
	var data graphql.Response
	state := utils.BuildDefaultState()
	newFish1 := generateFish("test1", "Test 1", "Desc 1")
	state.FishGateway.Fish[newFish1.Id] = newFish1
	newFish2 := generateFish("test2", "Test 2", "Desc 2")
	state.FishGateway.Fish[newFish2.Id] = newFish2

	utils.ExecuteGraphqlRequest(t, &state, query, "Fish", &data, nil)

	var fishListData FishResponse
	err := json.Unmarshal(data.Data, &fishListData)
	return state, data, fishListData, err
}

func generateFish(id string, name string, description string) *app.Fish {
	return &app.Fish{
		Id: id,
		Fish: &domain.Fish{
			Name:        name,
			Description: description,
		},
	}
}
