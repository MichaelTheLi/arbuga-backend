package main

import (
	"arbuga/backend/auth"
	"arbuga/backend/graph"
	"arbuga/backend/graph/model"
	"arbuga/backend/state"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"
const defaultFrontendUrl = "http://localhost:5173"
const defaultWebsocketOrigin = "localhost"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = defaultFrontendUrl
	}
	websocketOrigin := os.Getenv("WEBSOCKET_ORIGIN")
	if websocketOrigin == "" {
		websocketOrigin = defaultWebsocketOrigin
	}

	localState := state.AppLocalState{
		Users: make(map[string]*model.User),
	}
	router := chi.NewRouter()

	allowedHeaders := []string{"Authorization", "Content-Type"}
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{frontendUrl},
		AllowCredentials: true,
		AllowedHeaders:   allowedHeaders,
		//Debug:            true,
	}).Handler)

	router.Use(auth.Middleware(&localState))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{UsersState: &localState}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Host == websocketOrigin
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	// TODO Authentication https://gqlgen.com/recipes/authentication/
}
