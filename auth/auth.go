package auth

import (
	"arbuga/backend/graph/model"
	"arbuga/backend/state"
	"context"
	"encoding/base64"
	"net/http"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// Middleware AuthMiddleware decodes the share session cookie and packs the session into context
func Middleware(state *state.AppLocalState) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("auth-token")

			// Allow unauthenticated users in
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			userId, err := validateAndGetUserID(token)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			// get the user from the database
			user, _ := state.GetUserByID(userId)

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func validateAndGetUserID(token string) (string, error) {
	userId, err := base64.StdEncoding.DecodeString(token)
	return string(userId), err // TODO encode
}

// ForContext finds the user from the context. REQUIRES AuthMiddleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

// GenerateToken ForContext finds the user from the context. REQUIRES AuthMiddleware to have run.
func GenerateToken(user *model.User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.ID))
}
