package graph

import (
	"arbuga/backend/app"
	"arbuga/backend/domain"
	"context"
	"log"
	"net/http"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware AuthMiddleware decodes the share session cookie and packs the session into context
func Middleware(tokenService *app.TokenService, usersGateway *app.UserGateway) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const prefix = "Bearer "

			tokenHeader := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if tokenHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenValue := tokenHeader[len(prefix):]
			userId, err := (*tokenService).GetUserIdFromToken(tokenValue)

			// Checking token validity
			if err != nil {
				log.Println(err)
				log.Println(tokenValue)
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// get the user from the database
			user, _ := (*usersGateway).GetUserByID(userId)

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES AuthMiddleware to have run.
func ForContext(ctx context.Context) *domain.User {
	raw, _ := ctx.Value(userCtxKey).(*domain.User)
	return raw
}
