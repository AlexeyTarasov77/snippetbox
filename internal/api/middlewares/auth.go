package middlewares

import (
	"context"
	"net/http"

	"snippetbox.proj.net/internal/api/response"
)

// TODO: Cделать авторизацию по jwt токенам

func isValidToken(token string) bool {
	return token == "abc123"
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if token := r.Header.Get("Authorization"); token == "" || !isValidToken(token) {
				response.HttpError(w, "", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", "someuser")
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}