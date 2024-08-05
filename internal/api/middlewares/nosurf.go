package middlewares

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func NoSurfMiddleware(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}