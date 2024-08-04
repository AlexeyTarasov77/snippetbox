package models

import "time"

const SnippetQuery = `
		CREATE TABLE IF NOT EXISTS snippets (
			id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
			title varchar(100) NOT NULL,
			content text NOT NULL,
			created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires datetime NOT NULL
		)
	`

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}