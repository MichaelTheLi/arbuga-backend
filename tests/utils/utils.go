package utils

import (
	"arbuga/backend/api"
	"arbuga/backend/api/graph/model"
	"arbuga/backend/auth"
	"arbuga/backend/state"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BuildDefaultState() state.AppLocalState {
	return state.AppLocalState{
		Users: make(map[string]*model.User),
	}
}

func BuildStateWithUser(loginString string, passwordString string) state.AppLocalState {
	stateRes := BuildDefaultState()

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.MinCost) // TODO Reuse logic
	password := string(hashedPass)
	user := model.User{
		ID:       "testId",
		Login:    &loginString,
		Password: &password,
		Name:     "Test name",
	}
	stateRes.Users["testId"] = &user

	return stateRes
}

func ExecuteGraphqlRequest(t *testing.T, localState *state.AppLocalState, query string, operationName string, data any, token *string) {
	if localState == nil {
		defaultState := BuildDefaultState()
		localState = &defaultState
	}
	config := graph.BuildConfigFromEnv()

	queryEncoded, _ := json.Marshal(query)
	request := fmt.Sprintf("{\"query\":%s,\"operationName\":\"%s\"}", queryEncoded, operationName)
	body := strings.NewReader(request)

	req, err := http.NewRequest("POST", "/query", body)
	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	if token != nil {
		req.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", *token)}
	}

	assert.Nil(t, err, "Request created with an error")

	rr := httptest.NewRecorder()
	middleware := auth.Middleware(localState)
	handler := middleware(graph.BuildGraphqlServer(localState, config))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status != 200")
	jsonErr := json.NewDecoder(rr.Body).Decode(&data)
	assert.Nil(t, jsonErr, "Json not decoded")
}
