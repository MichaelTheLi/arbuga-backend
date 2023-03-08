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

type FishResponse struct {
	Fish []*model.Fish `json:"fish"`
}

func TestFishRequestWillReceiveData(t *testing.T) {
	state, _, fishListData, err := executeGetFish(t)
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, len(state.FishGateway.Fish))
}

func TestFishRequestWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "Bolivian ram")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, 1)
	assert.Equal(t, fishListData.Fish[0].Name, "Bolivian ram")
}

func TestFishRequestWillFindCorrectListOfFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "tetra")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, 2)

	assert.Equal(t, fishListData.Fish[0].Name, "Neon tetra")
	assert.Equal(t, fishListData.Fish[1].Name, "Rummy-nose tetra")
}

func TestFishRequestWithDifferentCaseWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "BoLiViAn RaM")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, 1)
	assert.Equal(t, fishListData.Fish[0].Name, "Bolivian ram")
}

func TestFishRequestWithRandomStringWillNotFindAnyFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "lkasfaglas asdas123 123 124ljkasdl")
	assert.Nil(t, err)

	assert.NotNil(t, fishListData.Fish)
	assert.Len(t, fishListData.Fish, 0)
}

func executeGetFish(t *testing.T) (utils.TestServerState, graphql.Response, FishResponse, error) {
	query := "query Fish {fish {id name description}}"
	var data graphql.Response
	state := BuildStateWithFish()

	utils.ExecuteGraphqlRequest(t, &state, query, "Fish", &data, nil)

	var fishListData FishResponse
	err := json.Unmarshal(data.Data, &fishListData)
	return state, data, fishListData, err
}

func executeSearchFish(t *testing.T, substring string) (utils.TestServerState, graphql.Response, FishResponse, error) {
	query := "query Fish($substring: String!) {fish(substring: $substring) {id name description}}"
	variables := fmt.Sprintf("{ \"substring\": \"%s\"}", substring)
	var data graphql.Response
	state := BuildStateWithFish()

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "Fish", &data, nil)

	var fishListData FishResponse
	err := json.Unmarshal(data.Data, &fishListData)
	return state, data, fishListData, err
}

func BuildStateWithFish() utils.TestServerState {
	state := utils.BuildDefaultState()
	newFish1 := generateFish("test1", "Rasbora", "Desc 1")
	state.FishGateway.Fish[newFish1.Id] = newFish1
	newFish2 := generateFish("test2", "Corydoras", "Desc 2")
	state.FishGateway.Fish[newFish2.Id] = newFish2
	newFish3 := generateFish("test3", "Bolivian ram", "Desc 3")
	state.FishGateway.Fish[newFish3.Id] = newFish3
	newFish4 := generateFish("test4", "Neon tetra", "Desc 4")
	state.FishGateway.Fish[newFish4.Id] = newFish4
	newFish5 := generateFish("test5", "Rummy-nose tetra", "Desc 5")
	state.FishGateway.Fish[newFish5.Id] = newFish5
	return state
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
