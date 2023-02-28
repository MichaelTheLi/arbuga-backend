package main

import (
	"arbuga/backend/adapters"
	"arbuga/backend/api"
	"arbuga/backend/api/graph"
	"arbuga/backend/app"
	"log"
	"net/http"
)

func main() {
	config := api.BuildConfigFromEnv()

	router := api.BuildServer(buildServerState(), config)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

// TODO ioc container?
func buildServerState() api.ServerState {
	userGateway := &adapters.LocalUserGateway{
		Users: make(map[string]*app.User),
	}
	tokenService := &adapters.JwtTokenService{
		Secret: "get_this_from_env", // TODO Get secret from env
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

	return api.ServerState{
		Resolver:     &resolver,
		TokenService: tokenService,
		UserGateway:  userGateway,
	}
}
