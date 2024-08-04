package storage

import "snippetbox.proj.net/internal/storage/models"

type SnippetsStorage interface {
	Insert(title, content string, expires int) (int64, error)
	Get(id int) (*models.Snippet, error)
	Latest(n int) ([]*models.Snippet, error)
}

type UsersStorage interface {
	Insert(username, email, password string) (int64, error)
	Authenticate(email, password string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Get(id int) (*models.User, error)
}