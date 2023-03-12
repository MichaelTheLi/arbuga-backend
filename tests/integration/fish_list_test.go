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

type FishListResponse struct {
	FishList []*model.Fish `json:"fishList"`
}

func TestFishRequestWillReceiveData(t *testing.T) {
	state, _, fishListData, err := executeGetFishList(t)
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.FishList)
	assert.Len(t, fishListData.FishList, len(state.FishGateway.Fish))
}

func TestFishRequestWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "Bolivian ram")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.FishList)
	assert.Len(t, fishListData.FishList, 1)
	assert.Equal(t, fishListData.FishList[0].Name, "Bolivian ram")
}

func TestFishRequestWillFindCorrectListOfFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "tetra")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.FishList)
	assert.Len(t, fishListData.FishList, 2)

	assert.Equal(t, fishListData.FishList[0].Name, "Neon tetra")
	assert.Equal(t, fishListData.FishList[1].Name, "Rummy-nose tetra")
}

func TestFishRequestWithDifferentCaseWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "BoLiViAn RaM")
	assert.Nil(t, err)

	assert.NotEmpty(t, fishListData.FishList)
	assert.Len(t, fishListData.FishList, 1)
	assert.Equal(t, fishListData.FishList[0].Name, "Bolivian ram")
}

func TestFishRequestWithRandomStringWillNotFindAnyFish(t *testing.T) {
	_, _, fishListData, err := executeSearchFish(t, "lkasfaglas asdas123 123 124ljkasdl")
	assert.Nil(t, err)

	assert.NotNil(t, fishListData.FishList)
	assert.Len(t, fishListData.FishList, 0)
}

func executeGetFishList(t *testing.T) (utils.TestServerState, graphql.Response, FishListResponse, error) {
	query := "query FishList {fishList {id name description}}"
	var data graphql.Response
	state := utils.BuildStateWithFish()

	utils.ExecuteGraphqlRequest(t, &state, query, "FishList", &data, nil)

	var fishListData FishListResponse
	err := json.Unmarshal(data.Data, &fishListData)
	return state, data, fishListData, err
}

func executeSearchFish(t *testing.T, substring string) (utils.TestServerState, graphql.Response, FishListResponse, error) {
	query := "query FishList($substring: String!) {fishList(substring: $substring) {id name description}}"
	variables := fmt.Sprintf("{ \"substring\": \"%s\"}", substring)
	var data graphql.Response
	state := utils.BuildStateWithFish()

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "FishList", &data, nil)

	var fishListData FishListResponse
	err := json.Unmarshal(data.Data, &fishListData)
	return state, data, fishListData, err
}
