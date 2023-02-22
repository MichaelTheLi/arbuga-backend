package main_test

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/auth"
	"arbuga/backend/tests/utils"
	json "encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type MeResponse struct {
	Me model.User `json:"me"`
}

func TestAuthenticatedWillReceiveData(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	state := utils.BuildDefaultState()
	login := "testLogin"
	passwordString := "testPass"

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.MinCost) // TODO Reuse logic
	password := string(hashedPass)
	user := model.User{
		ID:       "testId",
		Login:    &login,
		Password: &password,
		Name:     "Test name",
	}
	state.Users["testId"] = &user

	token, _ := auth.GenerateToken(&user)
	utils.ExecuteGraphqlRequest(t, &state, query, "Me", &data, &token)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	if err != nil {
		t.Errorf("Got err: %e", err)
	}

	if meData.Me.ID == "" {
		t.Errorf("Expected: %s. Got: %s.", "", meData.Me.ID)
	}
}

func TestNotAuthenticatedWillNotReceiveData(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	var meData MeResponse
	err := json.Unmarshal(data.Data, &meData)
	if err != nil {
		t.Errorf("Got err: %e", err)
	}

	if meData.Me.ID != "" {
		t.Errorf("Expected: %s. Got: %s.", "", meData.Me.ID)
	}
}

func TestNotAuthenticatedWillReceiveError(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	err := data.Errors[0]
	if err.Path.String() != "me" {
		t.Errorf("Expected: %s. Got: %s.", "me", err.Path.String())
	}

	if err.Message != "not authenticated" {
		t.Errorf("Expected: %s. Got: %s.", "not authenticated", err.Message)
	}
}

func TestNotAuthenticatedWillReceiveOnlyOneError(t *testing.T) {
	query := "query Me {me {id login name}}"
	var data graphql.Response
	utils.ExecuteGraphqlRequest(t, nil, query, "Me", &data, nil)

	if len(data.Errors) != 1 {
		t.Errorf("Expected 1 error, found %d", len(data.Errors))
	}
}
