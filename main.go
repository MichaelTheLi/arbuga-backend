package main

import (
	"arbuga/backend/api"
	"arbuga/backend/api/graph/model"
	"arbuga/backend/state"
	"log"
	"net/http"
)

func main() {
	localState := state.AppLocalState{
		Users: make(map[string]*model.User),
	}
	config := graph.BuildConfigFromEnv()

	router := graph.BuildServer(&localState, config)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
