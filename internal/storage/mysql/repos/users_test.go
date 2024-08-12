package repos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/models"
	"snippetbox.proj.net/internal/tests/utils"
)

func TestInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("repos: skipping test in short mode")
	}
	// t.Parallel()
	db := utils.NewTestDB(t)
	userModel := UserModel{db}
	dummyUser := utils.GetDummyUser()
	if _, err := db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		dummyUser.Username,
		dummyUser.Email,
		string(dummyUser.Password),
	); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name string
		username string
		email string
		password string
		expectedErr error
	} {
		{
			name: "valid user",
			username: "foo",
			email: "7JpjT@example.com",
			password: "qwerty123",
		},
		{
			name: "already existent user (duplicate email)",
			username: "foo",
			email: dummyUser.Email,
			password: "qwerty123",
			expectedErr: storage.ErrDuplicateEmail,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := userModel.Insert(tc.username, tc.email, tc.password)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Zero(t, id)
				return
			}
			assert.NoError(t, err)
			assert.NotZero(t, id)
			var insertedUser models.User
			if err := db.QueryRow(
				"SELECT username, email, password, is_active FROM users WHERE id = ?",
				id,
			).Scan(&insertedUser.Username, &insertedUser.Email, &insertedUser.Password, &insertedUser.IsActive);
			err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.username, insertedUser.Username)
			assert.Equal(t, tc.email, insertedUser.Email)
			assert.NotEqual(t, tc.password, string(insertedUser.Password))
		})
	}
}

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("repos: skipping test in short mode")
	}
	// t.Parallel()
	db := utils.NewTestDB(t)
	userModel := UserModel{db}
	dummyUser := utils.GetDummyUser()
	_, err := db.Exec(
		"INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)",
		dummyUser.ID,
		dummyUser.Username,
		dummyUser.Email,
		string(dummyUser.Password),
	)
	require.NoError(t, err)
	testCases := []struct {
		name string
		id int
		expectedErr error
	} {
		{
			name: "valid user",
			id: dummyUser.ID,
		},
		{
			name: "not found user",
			id: 0,
			expectedErr: storage.ErrNoRecord,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userModel.Get(tc.id)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, user)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, dummyUser.ID, user.ID)
			assert.Equal(t, dummyUser.Username, user.Username)
			assert.Equal(t, dummyUser.Email, user.Email)
		})
	}
}


func TestGetByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("repos: skipping test in short mode")
	}
	// t.Parallel()
	db := utils.NewTestDB(t)
	userModel := UserModel{db}
	dummyUser := utils.GetDummyUser()
	_, err := db.Exec(
		"INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)",
		dummyUser.ID,
		dummyUser.Username,
		dummyUser.Email,
		string(dummyUser.Password),
	)
	require.NoError(t, err)
	testCases := []struct {
		name string
		email string
		expectedErr error
	} {
		{
			name: "valid user",
			email: dummyUser.Email,
		},
		{
			name: "not found user",
			email: "7JpjT@example.com",
			expectedErr: storage.ErrNoRecord,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userModel.GetByEmail(tc.email)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, user)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, dummyUser.ID, user.ID)
			assert.Equal(t, dummyUser.Username, user.Username)
			assert.Equal(t, dummyUser.Email, user.Email)
		})
	}
}