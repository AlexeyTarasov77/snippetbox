package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"snippetbox.proj.net/internal/api/constants"
	"snippetbox.proj.net/internal/storage/models"
)

type UserGettableByID interface {
	Get(int) (*models.User, error)
}

func AuthMiddleware(logger *slog.Logger, sessionManager *scs.SessionManager, storage UserGettableByID) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				userId := sessionManager.GetInt(r.Context(), constants.UserIDCtxKey)
				var user *models.User
				user, err := storage.Get(userId)
				if err != nil {
					user = nil
				}
				r = r.WithContext(context.WithValue(r.Context(), constants.UserCtxKey, user))
				next.ServeHTTP(w, r)
			},
		)
	}
}