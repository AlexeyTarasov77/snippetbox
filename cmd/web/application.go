package main

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/schema"
	"snippetbox.proj.net/internal/storage"
)

type Application struct {
	logger        *slog.Logger
	snippets      storage.SnippetsStorage
	users         storage.UsersStorage
	templateCache map[string]*template.Template
	formDecoder       *schema.Decoder
	sessionManager *scs.SessionManager
}

func NewApplication(
	logger *slog.Logger, 
	snippets storage.SnippetsStorage,
	users storage.UsersStorage,
	templateCache map[string]*template.Template, 
	sessionManager *scs.SessionManager,
) *Application {
	return &Application{
		logger:        logger,
		users:         users,
		snippets:      snippets,
		templateCache: templateCache,
		formDecoder:   schema.NewDecoder(),
		sessionManager: sessionManager,
	}
}
