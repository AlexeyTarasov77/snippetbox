package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"snippetbox.proj.net/internal/tests/utils"
	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/mocks"
)

func TestPing(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)
	server := utils.NewTestServer(t, app.routes())
	defer server.Close()
	resp := server.Get("/ping")
	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "OK", resp.Body)
	// recorder := httptest.NewRecorder()
	// r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	// ping(recorder, r)
	// helpers.Equal(t, recorder.Code, http.StatusOK)
	// helpers.Equal(t, recorder.Body.String(), "OK")
}

// func TestSnippetView(t *testing.T) {
// 	t.Parallel()
// 	app := NewTestApplication(t)
// 	server := utils.NewTestServer(t, app.routes())
// 	defer server.Close()
// 	snippetsMock := app.snippets.(*mocks.SnippetsStorage)
// 	dummySnippet := utils.GetDummySnippet()
// 	snippetsMock.On("Get", 0).Return(nil, storage.ErrNoRecord)
// 	snippetsMock.On("Get", 1).Return(dummySnippet, nil)
// 	snippetsMock.On("Get", 2).Return(nil, storage.ErrNoRecord)
// 	testCases := []struct {
// 		name           string
// 		url            string
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name:           "Valid ID",
// 			url:            "/snippet/view/1",
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   dummySnippet.Content,
// 		},
// 		{
// 			name:           "Non-existent ID",
// 			url:            "/snippet/view/2",
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody:   http.StatusText(http.StatusNotFound),
// 		},
// 		{
// 			name:           "Negative ID",
// 			url:            "/snippet/view/-1",
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:           "Decimal ID",
// 			url:            "/snippet/view/1.23",
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:           "Empty ID",
// 			url:            "/snippet/view/",
// 			expectedStatus: http.StatusNotFound,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resp := server.Get(tc.url)
// 			assert.Equal(t, tc.expectedStatus, resp.Status)
// 			if tc.expectedBody != "" {
// 				assert.Contains(t, resp.Body, tc.expectedBody)
// 			}
// 		})
// 	}
// }

func TestUserSignup(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)
	usersMock := app.users.(*mocks.UsersStorage)
	dummyUser := utils.GetDummyUser()
	usersMock.On("Get", mock.AnythingOfType("int")).Return(dummyUser, nil)
	usersMock.On("Get", 0).Return(nil, storage.ErrNoRecord)
	usersMock.On(
		"Insert",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(int64(1), nil).Once()

	server := utils.NewTestServer(t, app.routes())
	defer server.Close()
	resp := server.Get("/user/signup")
	assert.Equal(t, resp.Status, http.StatusOK)
	validCSRFToken := utils.ExtractCSRFToken(t, resp.Body)
	t.Logf("CSRF token: %s", validCSRFToken)
	testCases := []struct {
		name             string
		username         string
		email            string
		password         string
		password_confirm string
		csrfToken        string
		expectedStatus   int
	}{
		{
			name:             "Valid submission",
			username:         dummyUser.Username,
			email:            dummyUser.Email,
			password:         string(dummyUser.Password),
			password_confirm: string(dummyUser.Password),
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusSeeOther,
		},
		{
			name:             "Invalid CSRF Token",
			username:         dummyUser.Username,
			email:            dummyUser.Email,
			password:         string(dummyUser.Password),
			password_confirm: string(dummyUser.Password),
			csrfToken:        "invalid",
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "Empty name",
			username:         "",
			email:            dummyUser.Email,
			password:         string(dummyUser.Password),
			password_confirm: string(dummyUser.Password),
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
		{
			name:             "Empty email",
			username:         dummyUser.Username,
			email:            "",
			password:         string(dummyUser.Password),
			password_confirm: string(dummyUser.Password),
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
		{
			name:             "Empty password",
			username:         dummyUser.Username,
			email:            dummyUser.Email,
			password:         "",
			password_confirm: string(dummyUser.Password),
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
		{
			name:             "Passwords don't match",
			username:         dummyUser.Username,
			email:            dummyUser.Email,
			password:         string(dummyUser.Password),
			password_confirm: "foobar123",
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
		{
			name:             "Password too short",
			username:         dummyUser.Username,
			email:            dummyUser.Email,
			password:         "foo",
			password_confirm: "foo",
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
		{
			name:             "Invalid email",
			username:         dummyUser.Username,
			email:            "foo",
			password:         string(dummyUser.Password),
			password_confirm: string(dummyUser.Password),
			csrfToken:        validCSRFToken,
			expectedStatus:   http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tc.username)
			form.Add("email", tc.email)
			form.Add("password", tc.password)
			form.Add("password_confirm", tc.password_confirm)
			form.Add("csrf_token", tc.csrfToken)
			resp = server.PostForm("/user/signup", form)
			assert.Equal(t, tc.expectedStatus, resp.Status)
		})
	}
}
