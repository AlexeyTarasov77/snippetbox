package repos

import (
	"database/sql"
	"errors"

	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/models"
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

func (model *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := model.DB.QueryRow(stmt, id)
	var snippet models.Snippet
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Snippet{}, storage.ErrNoRecord
		}
		return &models.Snippet{}, err
	}
	return &snippet, nil
}

func (model *SnippetModel) Latest(n int) ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT ?`
	rows, err := model.DB.Query(stmt, n)
	snippets := []*models.Snippet{}
	if err != nil {
		return snippets, err
	}
	defer rows.Close()
	for rows.Next() {
		var snippet models.Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return []*models.Snippet{}, err
		}
		snippets = append(snippets, &snippet)
	}
	if err = rows.Err(); err != nil {
		return []*models.Snippet{}, err
	}
	return snippets, nil
}