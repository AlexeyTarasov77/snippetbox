package middlewares

import (
	"log/slog"
	"net/http"

	"snippetbox.proj.net/internal/api/response"
)


func Recoverer(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func () {
					if err := recover(); err != nil {
						logger.Error("Recovered from panic", "msg", err)
						response.HttpError(w, "")
					}
				}()
				next.ServeHTTP(w, r)
			},
		)
	}
}