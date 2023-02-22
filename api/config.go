package graph

import (
	"os"
)

const defaultPort = "8080"
const defaultFrontendUrl = "http://localhost:5173"
const defaultWebsocketOrigin = "localhost"

type ServerConfig struct {
	Port                string
	CorsAllowedOrigins  []string
	CorsHeaders         []string
	CorsWebsocketOrigin string
}

func BuildConfigFromEnv() ServerConfig {
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

	return ServerConfig{
		Port:                port,
		CorsAllowedOrigins:  []string{frontendUrl},
		CorsHeaders:         []string{"Authorization", "Content-Type"},
		CorsWebsocketOrigin: websocketOrigin,
	}
}
