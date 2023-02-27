package api

import (
	"arbuga/backend/api/graph"
	"arbuga/backend/app"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type ServerState struct {
	Resolver     *graph.Resolver
	TokenService app.TokenService
	UserGateway  app.UserGateway
}

func BuildServer(serverState ServerState, config ServerConfig) *chi.Mux {
	router := buildRouter(serverState, config)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", BuildGraphqlServer(serverState.Resolver, config))

	return router
}

func BuildGraphqlServer(resolver *graph.Resolver, config ServerConfig) *handler.Server {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: resolver},
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

func buildRouter(serverState ServerState, config ServerConfig) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   config.CorsAllowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   config.CorsHeaders,
		//Debug:            true,
	}).Handler)

	router.Use(graph.Middleware(&serverState.TokenService, &serverState.UserGateway))

	return router
}
