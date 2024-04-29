package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/rs/zerolog/log"
)

const authorizationHeader = "Authorization"

type ctxUserID struct{}

func authenticate(auth auth.UserAuthentication) func(next http.Handler) http.Handler {
	// Подсмотрено в https://github.com/go-chi/chi/blob/master/middleware/content_type.go
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get(authorizationHeader)
			if ah == "" {
				log.Printf("Authorization header not set")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := auth.GetUserID(ah)
			if err != nil {
				log.Error().Err(err).Msg("error of get userID")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID{}, userID)
			log.Printf("Authentication UserID:%v", userID)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func getUserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(ctxUserID{}).(int64)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
