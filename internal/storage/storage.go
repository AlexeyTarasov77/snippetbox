package storage

import "time"

type ModelInterface interface {
	Insert(title, content string, expires int) (int64, error)
	Get(id int) (*Snippet, error)
	Latest(n int) ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
