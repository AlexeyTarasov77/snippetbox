package utils

import (
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"snippetbox.proj.net/internal/storage/models"
)

type testServer struct {
	*testing.T
	*httptest.Server
}

type TestResponse struct {
	Status int
	Body   string
	Headers http.Header
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

func readBody(t *testing.T, body io.ReadCloser) string {
	t.Helper()
	bodyBytes, err := io.ReadAll(body)
	defer body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(bodyBytes)
}

func (ts *testServer) Get(url string) TestResponse {
	resp, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		ts.Logf("ATTENTION: %s", err)
		ts.Fatal(err)
	}
	return TestResponse{resp.StatusCode, readBody(ts.T, resp.Body), resp.Header}
}

func (ts *testServer) PostForm(url string, data url.Values) TestResponse {
	resp, err := ts.Client().PostForm(ts.URL + url, data)
	if err != nil {
		ts.Fatal(err)
	}
	return TestResponse{resp.StatusCode, readBody(ts.T, resp.Body), resp.Header}
}

var csrfTokenRegexp = regexp.MustCompile(`<input\s+type\s*=\s*['"]hidden['"]\s+name\s*=\s*['"]csrf_token['"]\s+value\s*=\s*['"]([^'"]+)['"]`)

func ExtractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRegexp.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no CSRF token found in body")
	}
	return html.UnescapeString(matches[1])
}

func GetDummyUser() *models.User {
	return &models.User{
		ID:       rand.Int() + 1,
		Username: "foo",
		Email:    "NwJt2@example.com",
		Password: []byte("Pa$$w0rd"),
		Created:  time.Now(),
		IsActive: true,
	}
}

func GetDummySnippet() *models.Snippet {
	return &models.Snippet{
		ID:      rand.Int() + 1,
		Title:   "foo",
		Content: "bar",
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour * 24),
	}
}