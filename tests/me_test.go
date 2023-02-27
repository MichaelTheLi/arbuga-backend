package main_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/tests/utils"
	json "encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MeResponse struct {
	Me model.User `json:"me"`
}

func TestAuthenticatedWillReceiveData(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &state.Token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.NotEmpty(t, meData.Me.ID)
}

func TestNotAuthenticatedWillNotReceiveData(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	assert.Nil(t, err)

	assert.Empty(t, meData.Me.ID)
}

func TestNotAuthenticatedWillReceiveError(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	err := data.Errors[0]

	assert.Equal(t, "me", err.Path.String())
	assert.Equal(t, "not authenticated", err.Message)
}

func TestNotAuthenticatedWillReceiveOnlyOneError(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	assert.Len(t, data.Errors, 1)
}
