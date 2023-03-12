package main

import (
	"arbuga/backend/adapters"
	"arbuga/backend/api"
	"arbuga/backend/api/graph"
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"log"
	"net/http"
)

func main() {
	config := api.BuildConfigFromEnv()

	serverState := BuildServerState()
	seedServerState(&serverState)
	router := api.BuildServer(serverState, config)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func seedServerState(state *api.ServerState) {
	res, _ := state.Resolver.SignUpService.SignUp("test", "test", "Test User")
	width := 20
	height := 30
	length := 40
	thickness := 6
	name := "Test Ecosystem"
	state.Resolver.UserService.SaveEcosystem(res.User, &app.EcosystemInput{
		Name: &name,
		Aquarium: &app.AquariumGlassInput{
			Dimensions: &app.DimensionsInput{
				Width:  &width,
				Height: &height,
				Length: &length,
			},
			GlassThickness:     &thickness,
			SubstrateThickness: nil,
			DecorationsVolume:  nil,
		},
	})

	state.Resolver.FishService.Gateway.SaveFish(&app.Fish{
		Id: "test1",
		Fish: &domain.Fish{
			Name:        "Bolivian Ram",
			Description: "Just Ram",
		},
	})
	state.Resolver.FishService.Gateway.SaveFish(&app.Fish{
		Id: "test2",
		Fish: &domain.Fish{
			Name:        "Neon Tetra",
			Description: "Blue fish, fast one",
		},
	})
	state.Resolver.FishService.Gateway.SaveFish(&app.Fish{
		Id: "test3",
		Fish: &domain.Fish{
			Name:        "Harlequin Rasbora",
			Description: "With fancy pants",
		},
	})
	state.Resolver.FishService.Gateway.SaveFish(&app.Fish{
		Id: "test4",
		Fish: &domain.Fish{
			Name:        "Rummy-nose tetra",
			Description: "Red nose",
		},
	})
}

// TODO ioc container?
func BuildServerState() api.ServerState {
	userGateway := &adapters.LocalUserGateway{
		Users: make(map[string]*app.User),
	}
	fishGateway := &adapters.LocalFishGateway{
		Fish: make(map[string]*app.Fish),
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
	fishService := &app.FishService{
		Gateway: fishGateway,
	}
	resolver := graph.Resolver{
		SignInService: signInService,
		SignUpService: signUpService,
		UserService:   userService,
		FishService:   fishService,
	}

	return api.ServerState{
		Resolver:     &resolver,
		TokenService: tokenService,
		UserGateway:  userGateway,
	}
}
