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

func BuildStateWithUser(loginString string, passwordString string) (state.AppLocalState, string) {
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

	token, _ := auth.GenerateToken(stateRes.Users["testId"])
	return stateRes, token
}

func ExecuteGraphqlRequest(t *testing.T, localState *state.AppLocalState, query string, operationName string, data any, token *string) {
	executeGraphqlRequest(t, localState, query, "", operationName, data, token)
}

func ExecuteGraphqlRequestWithVariables(t *testing.T, localState *state.AppLocalState, query string, variables string, operationName string, data any, token *string) {
	executeGraphqlRequest(t, localState, query, variables, operationName, data, token)
}

func executeGraphqlRequest(t *testing.T, localState *state.AppLocalState, query string, variables string, operationName string, data any, token *string) {
	if localState == nil {
		defaultState := BuildDefaultState()
		localState = &defaultState
	}
	config := graph.BuildConfigFromEnv()

	var request string
	queryEncoded, _ := json.Marshal(query)
	if variables != "" {
		request = fmt.Sprintf("{\"query\":%s,\"operationName\":\"%s\", \"variables\":%s}", queryEncoded, operationName, variables)
	} else {
		request = fmt.Sprintf("{\"query\":%s,\"operationName\":\"%s\"}", queryEncoded, operationName)
	}
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

	assert.Equalf(t, http.StatusOK, rr.Code, "Status != 200. Body: %s", rr.Body.String())
	jsonErr := json.NewDecoder(rr.Body).Decode(&data)
	assert.Nil(t, jsonErr, "Json not decoded")
}
