package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"snippetbox.proj.net/internal/api/middlewares"
)

func (app *Application) registerRoutes(router *chi.Mux) {
	router.NotFound(app.snippetNotFound)
	router.Use(middlewares.Recoverer, middlewares.SecureHeaders)
	router.Route("/", func(router chi.Router) {
		router.Use(app.sessionManager.LoadAndSave)
		router.Get("/", app.home)
		router.Route("/snippet", func(router chi.Router) {
			router.Post("/create", app.snippetCreatePost)
			router.Get("/create", app.snippetCreateGet)
			router.Get("/view/{id}", app.snippetView)
		})
	})
	router.Handle(
		"/static/*",
		http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))),
	)
}
