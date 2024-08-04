package repos

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/models"
)

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(username, email, password string) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	stmt := `INSERT INTO users (username, email, password) 
		VALUES (?, ?, ?)`
	res, err := model.DB.Exec(stmt, username, email, passwordHash)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (model *UserModel) Authenticate(email, password string) (*models.User, error) {
	user, err := model.GetByEmail(email)
	if err != nil {
		return nil, storage.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, storage.ErrInvalidCredentials
	}
	return user, nil
}

func (model *UserModel) GetByEmail(email string) (*models.User, error) {
	res := model.DB.QueryRow("SELECT id, username, email, password, created, is_active FROM users WHERE email = ?", email)
	var user models.User
	err := res.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (model *UserModel) Get(id int) (*models.User, error) {
	res := model.DB.QueryRow("SELECT id, username, email, password, created, is_active FROM users WHERE id = ?", id)
	var user models.User
	err := res.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}