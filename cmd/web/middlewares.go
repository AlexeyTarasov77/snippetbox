package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
	"snippetbox.proj.net/internal/api/constants"
	"snippetbox.proj.net/internal/api/response"
	"snippetbox.proj.net/internal/storage/models"
)

type UserGettableByID interface {
	Get(int) (*models.User, error)
}

func (app *Application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := app.sessionManager.GetInt(r.Context(), string(constants.UserIDCtxKey))
			var user *models.User
			user, err := app.users.Get(userId)
			if err != nil {
				app.logger.Warn("Error getting user", "err", err.Error())
				user = nil
			}
			r = r.WithContext(context.WithValue(r.Context(), constants.UserCtxKey, user))
			next.ServeHTTP(w, r)
		},
	)
}

func (app *Application) LoginRequiredMiddleware(_shouldRedirect ...bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				userId := app.sessionManager.GetInt(r.Context(), string(constants.UserIDCtxKey))
				shouldRedirect := true
				if len(_shouldRedirect) > 0 {
					shouldRedirect = _shouldRedirect[0]
				}
				if userId == 0 {
					if shouldRedirect {
						app.sessionManager.Put(
							r.Context(),
							string(constants.RedirectCtxKey),
							r.URL.String(),
						)
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

func (app *Application) NoSurfMiddleware(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *Application) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					app.logger.Error("Recovered from panic", "msg", err)
					trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
					app.logger.Error("Stack trace", "msg", trace)
					errMsg := ""
					if app.debug {
						errMsg = trace
					}
					response.HttpError(w, errMsg)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}

func (app *Application) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Security-Policy",
				"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
			)
			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")
			next.ServeHTTP(w, r)
		},
	)
}
