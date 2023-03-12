package utils

import (
	"arbuga/backend/adapters"
	"arbuga/backend/api"
	"arbuga/backend/api/graph"
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestServerState struct {
	State        api.ServerState
	Token        string
	UsersGateway *adapters.LocalUserGateway
	FishGateway  *adapters.LocalFishGateway
}

func BuildDefaultState() TestServerState {
	userGateway := &adapters.LocalUserGateway{
		Users: make(map[string]*app.User),
	}
	fishGateway := &adapters.LocalFishGateway{
		Fish: make(map[string]*app.Fish),
	}
	tokenService := &adapters.JwtTokenService{
		Secret: "tests_secret",
	}
	authService := &adapters.BcryptAuthService{}
	signInService := &app.SignInService{
		Gateway:      userGateway,
		AuthService:  authService,
		TokenService: tokenService,
	}
	signUpService := &app.SignUpService{
		Gateway:     userGateway,
		AuthService: authService,
	}
	userService := &app.UserService{
		Gateway:     userGateway,
		FishGateway: fishGateway,
	}
	fishService := &app.FishService{
		Gateway: fishGateway,
	}
	resolver := graph.Resolver{
		SignInService: signInService,
		SignUpService: signUpService,
		UserService:   userService,
		FishService:   fishService,
	}

	return TestServerState{
		State: api.ServerState{
			Resolver:     &resolver,
			TokenService: tokenService,
			UserGateway:  userGateway,
		},
		Token:        "",
		UsersGateway: userGateway,
		FishGateway:  fishGateway,
	}
}

func BuildStateWithUser(loginString string, passwordString string) TestServerState {
	state := BuildDefaultState()

	state.UsersGateway.Users["testId"] = GenerateTestUser(loginString, passwordString, state)
	token, _ := state.State.TokenService.GenerateToken(state.UsersGateway.Users["testId"])

	state.Token = token

	return state
}

func BuildStateWithFish(state *TestServerState) TestServerState {
	if state == nil {
		defaultState := BuildDefaultState()
		state = &defaultState
	}
	newFish1 := generateFish("test1", "Rasbora", "Desc 1")
	state.FishGateway.Fish[newFish1.Id] = newFish1
	newFish2 := generateFish("test2", "Corydoras", "Desc 2")
	state.FishGateway.Fish[newFish2.Id] = newFish2
	newFish3 := generateFish("test3", "Bolivian ram", "Desc 3")
	state.FishGateway.Fish[newFish3.Id] = newFish3
	newFish4 := generateFish("test4", "Neon tetra", "Desc 4")
	state.FishGateway.Fish[newFish4.Id] = newFish4
	newFish5 := generateFish("test5", "Rummy-nose tetra", "Desc 5")
	state.FishGateway.Fish[newFish5.Id] = newFish5
	return *state
}

func generateFish(id string, name string, description string) *app.Fish {
	return &app.Fish{
		Id: id,
		Fish: &domain.Fish{
			Name:        name,
			Description: description,
		},
	}
}

func GenerateTestUser(loginString string, passwordString string, state TestServerState) *app.User {
	hashedPass, _ := state.State.Resolver.SignUpService.AuthService.HashFromPassword(passwordString)
	owner := &domain.Owner{
		Name: "Test name",
	}
	user := &app.User{
		ID:           "testId",
		Login:        &loginString,
		PasswordHash: &hashedPass,
		Owner:        owner,
	}
	return user
}

func ExecuteGraphqlRequest(t *testing.T, serverState *TestServerState, query string, operationName string, data any, token *string) {
	executeGraphqlRequest(t, serverState, query, "", operationName, data, token)
}

func ExecuteGraphqlRequestWithVariables(t *testing.T, serverState *TestServerState, query string, variables string, operationName string, data any, token *string) {
	executeGraphqlRequest(t, serverState, query, variables, operationName, data, token)
}

func executeGraphqlRequest(t *testing.T, serverState *TestServerState, query string, variables string, operationName string, data any, token *string) {
	if serverState == nil {
		defaultState := BuildDefaultState()
		serverState = &defaultState
	}
	config := api.BuildConfigFromEnv()

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
	middleware := graph.Middleware(&serverState.State.TokenService, &serverState.State.UserGateway)
	handler := middleware(api.BuildGraphqlServer(serverState.State.Resolver, config))
	handler.ServeHTTP(rr, req)

	assert.Equalf(t, http.StatusOK, rr.Code, "Status != 200. Body: %s", rr.Body.String())
	jsonErr := json.NewDecoder(rr.Body).Decode(&data)
	assert.Nil(t, jsonErr, "Json not decoded")
}
