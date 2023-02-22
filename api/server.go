package graph

import (
	"arbuga/backend/api/graph"
	"arbuga/backend/auth"
	"arbuga/backend/state"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func BuildServer(localState *state.AppLocalState, config ServerConfig) *chi.Mux {
	router := buildRouter(localState, config)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", BuildGraphqlServer(localState, config))

	return router
}

func BuildGraphqlServer(localState *state.AppLocalState, config ServerConfig) *handler.Server {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{UsersState: localState}},
		),
	)

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Host == config.CorsWebsocketOrigin
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	return srv
}

func buildRouter(localState *state.AppLocalState, config ServerConfig) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   config.CorsAllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   config.CorsHeaders,
		//Debug:            true,
	}).Handler)

	router.Use(auth.Middleware(localState))

	return router
}
