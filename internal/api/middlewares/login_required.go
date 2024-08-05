package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"snippetbox.proj.net/internal/api/constants"
	"snippetbox.proj.net/internal/api/response"
)

func LoginRequiredMiddleware(
	logger *slog.Logger,
    sessionManager *scs.SessionManager,
	storage UserGettableByID, _shouldRedirect ...bool,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				userId := sessionManager.GetInt(r.Context(), string(constants.UserIDCtxKey))
				shouldRedirect := true
				if len(_shouldRedirect) > 0 {
					shouldRedirect = _shouldRedirect[0]
				}
				if userId == 0 {
					if shouldRedirect {
						http.Redirect(w, r, "/user/login", http.StatusSeeOther)
					} else {
						response.HttpError(w, "", http.StatusUnauthorized)
					}
                   return
				}
				w.Header().Add("Cache-Control", "no-store")
				next.ServeHTTP(w, r)
			},
		)
	}
}