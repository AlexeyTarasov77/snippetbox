package main

import (
	"net/http"
	"testing"

	"snippetbox.proj.net/internal/lib/tests/helpers"
	"snippetbox.proj.net/internal/tests/utils"
)

func TestPing(t *testing.T) {
	app := NewTestApplication()
	server := utils.NewTestServer(t, app.routes())
	defer server.Close()
	status, _, body := server.Get("/ping")
	helpers.Equal(t, status, http.StatusOK)
	helpers.Equal(t, string(body), "OK")
	// recorder := httptest.NewRecorder()
	// r := httptest.NewRequest(http.MethodGet, "/ping", nil)
	// ping(recorder, r)
	// helpers.Equal(t, recorder.Code, http.StatusOK)
	// helpers.Equal(t, recorder.Body.String(), "OK")
}
