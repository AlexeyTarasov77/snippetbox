package models

import "time"

const UserQuery = `CREATE TABLE IF NOT EXISTS users (
			id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
			username varchar(255) NOT NULL,
			email varchar(255) NOT NULL UNIQUE,
			password char(60) NOT NULL,
			created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_active boolean NOT NULL DEFAULT TRUE
		)`

type User struct {
	ID       int
	Username string
	Email    string
	Password []byte
	Created  time.Time
	IsActive   bool
}

