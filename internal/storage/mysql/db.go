package mysql

import (
	"database/sql"
	"log/slog"

	"snippetbox.proj.net/internal/storage/models"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		slog.Error("Error connecting to database", "msg", err)
		return nil, err
	} 
	if err := db.Ping(); err != nil {
		slog.Error("Error pinging database", "msg", err)
		return nil, err
	}
	slog.Info("Database connection established")
	_, err = db.Exec(models.SnippetQuery)
	if err != nil {
		slog.Error("Error creating snippets table", "msg", err)
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			token CHAR(43) PRIMARY KEY,
			data BLOB NOT NULL,
			expiry TIMESTAMP(6) NOT NULL
		);
	`)
	db.Exec(`CREATE INDEX sessions_expiry_idx ON sessions (expiry);`)
	if err != nil {
		slog.Error("Error creating sessions table", "msg", err)
		return nil, err
	}
	_, err = db.Exec(models.UserQuery)
	if err != nil {
		slog.Error("Error creating users table", "msg", err)
		return nil, err
	}
	slog.Info("Database tables created")
	return db, nil
}