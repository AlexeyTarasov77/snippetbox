package main

import (
	"flag"
	"fmt"

	"log"
	"log/slog"
	"net/http"
	"os"

	"snippetbox.proj.net/internal/config"
	"snippetbox.proj.net/internal/storage/mysql"
	"snippetbox.proj.net/internal/storage/mysql/models"
	"snippetbox.proj.net/internal/templates"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := flag.Int("port", 0, "The port of the application.")
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("Config path is required")
	}
	if _, err := os.Stat(*configPath); err != nil {
		log.Fatalf("Config file not found: %s", *configPath)
	}
	config, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if *port == 0 {
		*port = config.HTTPServer.Port
	}
	logger := setupLogger()
	serverAddr := fmt.Sprintf("%s:%d", config.HTTPServer.Host, *port)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name,
	)
	db, err := mysql.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Cookie.Secure = true
	app := NewApplication(logger, &models.SnippetModel{DB: db}, templateCache, sessionManager)
	defer db.Close()
	router := chi.NewRouter()
	app.registerRoutes(router)
	server := http.Server{Addr: serverAddr, Handler: router}
	slog.Info(fmt.Sprintf("Starting server on https://%s", serverAddr))
	if err := server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
		slog.Error(fmt.Sprintf("Server crashed with error %s", err))
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}