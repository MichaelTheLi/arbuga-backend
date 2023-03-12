package integration_test

import (
	"arbuga/backend/tests/integration/utils"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyUserDoesntHaveFish(t *testing.T) {
	_, _, meData := executeGetMe(t)
	assert.Empty(t, meData.Me.Fish)
}

func executeGetMe(t *testing.T) (utils.TestServerState, graphql.Response, MeResponse) {
	query := "query Me {me {id fish {id name}}}"
	var data graphql.Response
	stateWithUser := utils.BuildStateWithUser("testLogin", "testPass")
	state := utils.BuildStateWithFish(&stateWithUser)

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, nil)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)
	return state, data, meData
}
