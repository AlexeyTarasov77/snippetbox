package main

import (
	"html/template"
	"log/slog"
	"os"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/schema"
	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/mocks"
	"snippetbox.proj.net/internal/templates"
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
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return &Application{
		logger:        logger,
		users:         users,
		snippets:      snippets,
		templateCache: templateCache,
		formDecoder:   decoder,
		sessionManager: sessionManager,
	}
}

func NewTestApplication(t testing.TB) *Application {
	tc, err := templates.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	sessionManager := scs.New()
	sessionManager.Cookie.Secure = true
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return &Application{
		logger:  slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
		snippets: mocks.NewSnippetsStorage(t),
		users:   mocks.NewUsersStorage(t),
		templateCache: tc,
		formDecoder: decoder,
		sessionManager: sessionManager,
	}
}