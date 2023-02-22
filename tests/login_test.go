package main_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/tests/utils"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
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

	login, _ := state.GetUserByLogin("testLogin")
	if login == nil {
		t.Error("Expected the user to be created, actually - not created")
	}

	if len(state.Users) != 1 {
		t.Errorf("Expected only 1 user to be created, got %d", len(state.Users))
	}
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

	if len(state.Users) != 1 {
		if len(state.Users) > 1 {
			t.Error("Expected that existing user will be used, got new created")
		} else {
			t.Error("Expected that existing user will be used, got removed somehow")
		}
	}
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

	if err != nil {
		t.Errorf("Got err: %e", err)
	}

	if loginData.Login != nil {
		t.Error("Expected to not receive user, got one")
	}
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
	if err.Path.String() != "login" {
		t.Errorf("Expected: %s. Got: %s.", "login", err.Path.String())
	}

	if err.Message != "error" {
		t.Errorf("Expected: %s. Got: %s.", "error", err.Message)
	}
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

	if err != nil {
		t.Errorf("Received error %e", err)
	}
	if loginData.Login.Token == nil || *loginData.Login.Token == "" {
		t.Error("Expected to receive, got none")
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

	if err != nil {
		t.Errorf("Got err: %e", err)
	}

	if loginData.Login.User == nil {
		t.Error("Expected to receive user, got none")
	}

	if *loginData.Login.User.Login != "testLogin" {
		t.Errorf("Expected to receive the same Login, got %s", *loginData.Login.User.Login)
	}
}
