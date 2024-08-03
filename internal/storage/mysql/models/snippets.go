package models

import (
	"database/sql"
	"errors"

	"snippetbox.proj.net/internal/storage"
)

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int64, error) {
	stmt := `INSERT INTO snippets (title, content, expires) 
		VALUES (?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := model.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (model *SnippetModel) Get(id int) (*storage.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := model.DB.QueryRow(stmt, id)
	var snippet storage.Snippet
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &storage.Snippet{}, storage.ErrNoRecord
		}
		return &storage.Snippet{}, err
	}
	return &snippet, nil
}

func (model *SnippetModel) Latest(n int) ([]*storage.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT ?`
	rows, err := model.DB.Query(stmt, n)
	snippets := []*storage.Snippet{}
	if err != nil {
		return snippets, err
	}
	defer rows.Close()
	for rows.Next() {
		var snippet storage.Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return []*storage.Snippet{}, err
		}
		snippets = append(snippets, &snippet)
	}
	if err = rows.Err(); err != nil {
		return []*storage.Snippet{}, err
	}
	return snippets, nil
}