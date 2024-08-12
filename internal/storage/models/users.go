package models

import "time"


type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
	Created  time.Time
	IsActive   bool
}

