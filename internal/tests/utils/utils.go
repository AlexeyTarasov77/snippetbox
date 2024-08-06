package utils

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*testing.T
	*httptest.Server
}

func NewTestServer(t *testing.T, handler http.Handler) *testServer {
	ts := httptest.NewTLSServer(handler)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{t, ts}
}

func (ts *testServer) Get(url string) (status int, headers http.Header, body string) {
	resp, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		ts.Fatal(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ts.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode, resp.Header, string(bodyBytes)
}