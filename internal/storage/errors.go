package storage

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials provided")
	ErrDuplicateEmail = errors.New("models: user with that email already exists")
)