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
	Connection model.FishListConnection `json:"fishList"`
}

func TestFishRequestWillReceiveData(t *testing.T) {
	state, _, fishListData := executeGetFishList(t)

	assert.Len(t, fishListData.Connection.Edges, len(state.FishGateway.Fish))
}

func TestFishRequestWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData := executeSearchFish(t, "Bolivian ram")

	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, fishListData.Connection.Edges[0].Node.Name, "Bolivian ram")
}

func TestFishRequestWillFindCorrectListOfFish(t *testing.T) {
	_, _, fishListData := executeSearchFish(t, "tetra")

	assert.Len(t, fishListData.Connection.Edges, 3)

	assert.Equal(t, fishListData.Connection.Edges[0].Node.Name, "Neon tetra")
	assert.Equal(t, fishListData.Connection.Edges[1].Node.Name, "Rummy-nose tetra")
	assert.Equal(t, fishListData.Connection.Edges[2].Node.Name, "Yet another tetra")
}

func TestFishRequestWithDifferentCaseWillFindCorrectSingleFish(t *testing.T) {
	_, _, fishListData := executeSearchFish(t, "BoLiViAn RaM")

	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, fishListData.Connection.Edges[0].Node.Name, "Bolivian ram")
}

func TestFishRequestWithRandomStringWillNotFindAnyFish(t *testing.T) {
	_, data, fishListData := executeSearchFish(t, "lkasfaglas asdas123 123 124ljkasdl")

	assert.Len(t, fishListData.Connection.Edges, 0)
	assert.Nil(t, data.Errors)
}

func TestFishRequestPaginationCursorWorks(t *testing.T) {
	var _, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, nil)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, "Neon tetra", fishListData.Connection.Edges[0].Node.Name)

	_, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, &fishListData.Connection.PageInfo.EndCursor)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, "Rummy-nose tetra", fishListData.Connection.Edges[0].Node.Name)

	_, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, &fishListData.Connection.PageInfo.EndCursor)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, "Yet another tetra", fishListData.Connection.Edges[0].Node.Name)
}

func TestFishRequestPaginationHasNextPageWorks(t *testing.T) {
	var _, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, nil)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, true, *fishListData.Connection.PageInfo.HasNextPage)

	_, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, &fishListData.Connection.PageInfo.EndCursor)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, true, *fishListData.Connection.PageInfo.HasNextPage)

	_, _, fishListData = executeSearchFishPaginated(t, "Tetra", 1, &fishListData.Connection.PageInfo.EndCursor)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, false, *fishListData.Connection.PageInfo.HasNextPage)
}

func TestFishRequestPaginationHasNextPageForIncompletePageWorks(t *testing.T) {
	var _, _, fishListData = executeSearchFishPaginated(t, "Tetra", 2, nil)
	assert.Len(t, fishListData.Connection.Edges, 2)
	assert.Equal(t, true, *fishListData.Connection.PageInfo.HasNextPage)

	_, _, fishListData = executeSearchFishPaginated(t, "Tetra", 2, &fishListData.Connection.PageInfo.EndCursor)
	assert.Len(t, fishListData.Connection.Edges, 1)
	assert.Equal(t, false, *fishListData.Connection.PageInfo.HasNextPage)
}

func TestFishRequestPaginationHasNextPageForTooManyWorks(t *testing.T) {
	var _, _, fishListData = executeSearchFishPaginated(t, "Tetra", 10, nil)
	assert.Len(t, fishListData.Connection.Edges, 3)
	assert.Equal(t, false, *fishListData.Connection.PageInfo.HasNextPage)
}

func executeGetFishList(t *testing.T) (utils.TestServerState, graphql.Response, FishListResponse) {
	query := "query FishList {fishList { edges { cursor node { id name description }} pageInfo {startCursor endCursor hasNextPage}}}"
	var data graphql.Response
	state := utils.BuildStateWithFish(nil)

	utils.ExecuteGraphqlRequest(t, &state, query, "FishList", &data, nil)

	var fishListData FishListResponse
	err := json.Unmarshal(data.Data, &fishListData)
	assert.Nil(t, err)
	return state, data, fishListData
}

func executeSearchFish(t *testing.T, substring string) (utils.TestServerState, graphql.Response, FishListResponse) {
	query := "query FishList($substring: String!) {fishList(substring: $substring) { edges { cursor node { id name description }} pageInfo {startCursor endCursor hasNextPage}}}"
	variables := fmt.Sprintf("{ \"substring\": \"%s\"}", substring)
	var data graphql.Response
	state := utils.BuildStateWithFish(nil)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "FishList", &data, nil)

	var fishListData FishListResponse
	err := json.Unmarshal(data.Data, &fishListData)
	assert.Nil(t, err)
	return state, data, fishListData
}

func executeSearchFishPaginated(t *testing.T, substring string, first int, after *string) (utils.TestServerState, graphql.Response, FishListResponse) {
	query := "query FishList($substring: String!, $first: Int!, $after: ID) {fishList(substring: $substring, first: $first, after: $after) { edges { cursor node { id name description }} pageInfo {startCursor endCursor hasNextPage}}}"

	var variables string
	if after != nil {
		variables = fmt.Sprintf("{ \"substring\": \"%s\", \"first\": \"%d\", \"after\": \"%s\"}", substring, first, *after)
	} else {
		variables = fmt.Sprintf("{ \"substring\": \"%s\", \"first\": \"%d\"}", substring, first)
	}

	var data graphql.Response
	state := utils.BuildStateWithFish(nil)

	utils.ExecuteGraphqlRequestWithVariables(t, &state, query, variables, "FishList", &data, nil)

	var fishListData FishListResponse
	err := json.Unmarshal(data.Data, &fishListData)
	assert.Nil(t, err)
	return state, data, fishListData
}
