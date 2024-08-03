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
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS snippets (
			id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
			title varchar(100) NOT NULL,
			content text NOT NULL,
			created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires datetime NOT NULL
		)
	`)
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
	_, err2 := db.Exec(`CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);`)
	if err != nil{
		slog.Error("Error creating sessions table", "msg", err2)
		return nil, err
	}
	slog.Info("Database tables created")
	return db, nil
}