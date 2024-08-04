package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"snippetbox.proj.net/internal/api/middlewares"
)

func (app *Application) registerRoutes(router *chi.Mux) {
	router.NotFound(app.snippetNotFound)
	router.Use(middlewares.Recoverer(app.logger), middlewares.SecureHeaders)
	router.Route("/", func(router chi.Router) {
		router.Use(app.sessionManager.LoadAndSave, middlewares.AuthMiddleware(app.logger, app.sessionManager, app.users))
		router.Get("/", app.home)
		router.Route("/snippet", func(router chi.Router) {
			router.Post("/create", app.snippetCreatePost)
			router.Get("/create", app.snippetCreate)
			router.Get("/view/{id}", app.snippetView)
		})
		router.Route("/user", func(router chi.Router) {
			router.Get("/signup", app.userSignup)
			router.Post("/signup", app.userSignupPost)
			router.Get("/login", app.userLogin)
			router.Post("/login", app.userLoginPost)
			router.Post("/logout", app.userLogoutPost)
		})
	})
	router.Handle(
		"/static/*",
		http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))),
	)
}
