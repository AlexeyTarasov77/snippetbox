package mysql

import (
	"database/sql"
	"log/slog"
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
	return db, nil
}