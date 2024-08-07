package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSecureHeaders(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })
	app := NewTestApplication(t)
	app.SecureHeaders(next).ServeHTTP(recorder, r)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, recorder.Body.String(), "OK")
	assert.Equal(
		t, recorder.Header().Get("Content-Security-Policy"),
	 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
	)
	assert.Equal(t, recorder.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, recorder.Header().Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, recorder.Header().Get("X-Frame-Options"), "deny")
	assert.Equal(t, recorder.Header().Get("X-XSS-Protection"), "0")
}
