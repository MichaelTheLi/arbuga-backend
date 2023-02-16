package auth

import (
	"arbuga/backend/graph/model"
	"arbuga/backend/state"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
)

const secret = "t0k3n" // TODO Get from env

type UserClaims struct {
	jwt.RegisteredClaims
	UserId string
	Name   string
}

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
			const prefix = "Bearer "

			tokenHeader := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if tokenHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenValue := tokenHeader[len(prefix):]
			userClaims := &UserClaims{}
			token, err := jwt.ParseWithClaims(tokenValue, userClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			// Checking token validity
			if err != nil || !token.Valid {
				log.Println(err)
				log.Println(token)
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			userId := userClaims.UserId

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

// ForContext finds the user from the context. REQUIRES AuthMiddleware to have run.
func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}

func GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), UserClaims{
		UserId:           user.ID,
		Name:             user.Name,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
