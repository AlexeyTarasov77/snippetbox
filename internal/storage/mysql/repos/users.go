package repos

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
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
	var mySQLErr *mysql.MySQLError
	if err != nil {
		if errors.As(err, &mySQLErr) && mySQLErr.Number == 1062 {
			return 0, storage.ErrDuplicateEmail
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, storage.ErrInvalidCredentials
	}
	return user, nil
}

func (model *UserModel) getUser (query string, params ...any) (*models.User, error) {
	res := model.DB.QueryRow(query, params...)
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

func (model *UserModel) GetByEmail(email string) (*models.User, error) {
	return model.getUser(
		"SELECT id, username, email, password, created, is_active FROM users WHERE email = ? AND is_active = 1",
		email,
	)
}

func (model *UserModel) Get(id int) (*models.User, error) {
	return model.getUser(
		"SELECT id, username, email, password, created, is_active FROM users WHERE id = ? AND is_active = 1",
		id,
	)
}