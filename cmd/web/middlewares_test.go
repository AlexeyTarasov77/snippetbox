package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.proj.net/internal/lib/tests/helpers"
)


func TestSecureHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })
	app := NewTestApplication()
	app.SecureHeaders(next).ServeHTTP(recorder, r)
	helpers.Equal(t, recorder.Code, http.StatusOK)
	helpers.Equal(t, recorder.Body.String(), "OK")
	helpers.Equal(
		t, recorder.Header().Get("Content-Security-Policy"),
	 	"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
	)
	helpers.Equal(t, recorder.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
	helpers.Equal(t, recorder.Header().Get("X-Content-Type-Options"), "nosniff")
	helpers.Equal(t, recorder.Header().Get("X-Frame-Options"), "deny")
	helpers.Equal(t, recorder.Header().Get("X-XSS-Protection"), "0")
}
