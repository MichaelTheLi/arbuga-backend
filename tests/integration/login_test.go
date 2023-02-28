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

type LoginResponse struct {
	Login *model.LoginResult `json:"login"`
}

func TestLoginCreatesUser(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"testPass",
	)

	var data graphql.Response
	state := utils.BuildDefaultState()
	utils.ExecuteGraphqlRequest(t, &state, query, "LoginUser", &data, nil)

	user, err := state.State.UserGateway.GetUserByLogin("testLogin")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Len(t, state.UsersGateway.Users, 1)
}

func TestLoginLoginsExistingUser(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"testPass",
	)

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequest(t, &state, query, "LoginUser", &data, nil)

	assert.Len(t, state.UsersGateway.Users, 1, "Should be still single user")
}

func TestLoginWontLoginWithWrongPassword(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"wrongPassword",
	)

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequest(t, &state, query, "LoginUser", &data, nil)

	var loginData LoginResponse
	err := json.Unmarshal(data.Data, &loginData)
	assert.Nil(t, err)
	assert.Nil(t, loginData.Login)
}

func TestLoginWillReceiveErrorWithWrongPassword(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"wrongPassword",
	)

	var data graphql.Response
	state := utils.BuildStateWithUser("testLogin", "testPass")

	utils.ExecuteGraphqlRequest(t, &state, query, "LoginUser", &data, nil)

	err := data.Errors[0]

	assert.Equal(t, "login", err.Path.String())
	assert.Equal(t, "error: invalid password", err.Message)
}

func TestLoginReturnsTokenWithValidData(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"testPass",
	)

	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "LoginUser", &data, nil)

	var loginData LoginResponse
	err := json.Unmarshal(data.Data, &loginData)
	assert.Nil(t, err)
	if loginData.Login != nil {
		assert.NotEmpty(t, loginData.Login.Token)
	} else {
		assert.Fail(t, "Wrong response")
	}
}

func TestLoginReturnsUserWithValidData(t *testing.T) {
	query := fmt.Sprintf(
		"mutation LoginUser{ login(login: \"%s\", password: \"%s\") { token user { id login } } }",
		"testLogin",
		"testPass",
	)

	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "LoginUser", &data, nil)

	var loginData LoginResponse
	err := json.Unmarshal(data.Data, &loginData)
	assert.Nil(t, err)

	if loginData.Login != nil {
		assert.NotNil(t, loginData.Login.User)
		assert.Equal(t, "testLogin", *loginData.Login.User.Login)
	} else {
		assert.Fail(t, "Wrong response")
	}
}
