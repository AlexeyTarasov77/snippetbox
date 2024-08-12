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

func (model *SnippetModel) Insert(title, content string, expires int, userID int) (int64, error) {
	stmt := `INSERT INTO snippets (title, content, expires, user_id) 
		VALUES (?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`

	res, err := model.DB.Exec(stmt, title, content, expires, userID)
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
	stmt := `
		SELECT s.id, s.title, s.content, s.created, s.expires, u.username FROM snippets s
		JOIN users u ON s.user_id = u.id
		WHERE s.expires > UTC_TIMESTAMP() AND s.id = ?`
	row := model.DB.QueryRow(stmt, id)
	var snippet models.Snippet
	var user models.User
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoRecord
		}
		return nil, err
	}
	snippet.User = &user
	return &snippet, nil
}

func (model *SnippetModel) GetByUserID(userID int) ([]*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE user_id = ? ORDER BY created DESC
	`
	rows, err := model.DB.Query(stmt, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Snippet{}, storage.ErrNoRecord
		}
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
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

func (model *SnippetModel) Latest(n int) ([]*models.Snippet, error) {
	stmt := `SELECT s.id, s.title, s.content, s.created, s.expires, u.username FROM snippets s
			JOIN users u ON s.user_id = u.id
			WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT ?`
	rows, err := model.DB.Query(stmt, n)
	snippets := []*models.Snippet{}
	if err != nil {
		return snippets, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		var snippet models.Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires, &user.Username)
		if err != nil {
			return []*models.Snippet{}, err
		}
		snippet.User = &user
		snippets = append(snippets, &snippet)
	}
	// println(snippets[0].User.Username)
	if err = rows.Err(); err != nil {
		return []*models.Snippet{}, err
	}
	return snippets, nil
}