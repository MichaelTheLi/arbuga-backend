package integration_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FishResponse struct {
	Fish *model.Fish `json:"fish"`
}

func TestFishFoundById(t *testing.T) {
	idInTest := "test2"
	state, _, fishData := executeGetFish(t, idInTest)

	assert.NotEmpty(t, fishData.Fish)
	assert.Equal(t, fishData.Fish.ID, idInTest)
	assert.Equal(t, fishData.Fish.Name, state.FishGateway.Fish[idInTest].Fish.Name)
}

func TestFishNoFoundByInvalidId(t *testing.T) {
	idInTest := "mz1xa23jgl5sal"
	_, data, fishData := executeGetFish(t, idInTest)

	assert.Empty(t, fishData.Fish)
	assert.NotEmpty(t, data.Errors)
	assert.Len(t, data.Errors, 1)
	assert.Equal(t, "fish", data.Errors[0].Path.String())
	assert.Equal(t, "fish not found", data.Errors[0].Message)
}

func executeGetFish(t *testing.T, id string) (utils.TestServerState, graphql.Response, FishResponse) {
	query := "query Fish($id: ID!) {fish(id: $id) {id name description}}"
	variables := fmt.Sprintf("{ \"id\": \"%s\"}", id)
	var data graphql.Response
	state := utils.BuildStateWithFish(nil)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "Fish", &data, nil)

	var fishData FishResponse
	err := json.Unmarshal(data.Data, &fishData)
	assert.Nil(t, err)

	return state, data, fishData
}
