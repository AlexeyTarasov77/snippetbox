package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"snippetbox.proj.net/ui"
)

func (app *Application) routes() *chi.Mux {
	router := chi.NewRouter()

	router.NotFound(app.snippetNotFound)
	router.Use(app.Recoverer, app.SecureHeaders)
	router.Get("/ping", app.ping)
	router.Route("/", func(router chi.Router) {
		router.Use(app.sessionManager.LoadAndSave)
		router.Use(app.AuthMiddleware)
		router.Use(app.NoSurfMiddleware)
		router.Get("/", app.home)
		router.Route("/snippet", func(router chi.Router) {
			router.Route("/", func(router chi.Router) { // protected routes
				router.Use(app.LoginRequiredMiddleware())
				router.Post("/create", app.snippetCreatePost)
				router.Get("/create", app.snippetCreate)
			})
			router.Get("/view/{id}", app.snippetView)
		})
		router.Route("/user", func(router chi.Router) {
			router.Get("/signup", app.userSignup)
			router.Post("/signup", app.userSignupPost)
			router.Get("/login", app.userLogin)
			router.Post("/login", app.userLoginPost)
			router.With(
				app.LoginRequiredMiddleware(),
			).Post("/logout", app.userLogoutPost)
		})
	})
	router.Handle("/static/*", http.FileServer(http.FS(ui.Files)))

	return router
}
