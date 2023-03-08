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
		Gateway: userGateway,
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

// TODO Simplify
func BuildStateWithUser(loginString string, passwordString string) TestServerState {
	userGateway := &adapters.LocalUserGateway{
		Users: make(map[string]*app.User),
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
		Gateway: userGateway,
	}
	resolver := graph.Resolver{
		SignInService: signInService,
		SignUpService: signUpService,
		UserService:   userService,
	}

	hashedPass, _ := authService.HashFromPassword(passwordString)
	owner := &domain.Owner{
		Name: "Test name",
	}
	user := &app.User{
		ID:           "testId",
		Login:        &loginString,
		PasswordHash: &hashedPass,
		Owner:        owner,
	}
	userGateway.Users["testId"] = user

	token, _ := tokenService.GenerateToken(user)

	return TestServerState{
		State: api.ServerState{
			Resolver:     &resolver,
			TokenService: tokenService,
			UserGateway:  userGateway,
		},
		Token:        token,
		UsersGateway: userGateway,
	}
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
