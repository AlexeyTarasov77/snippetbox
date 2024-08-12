package utils

import (
	"database/sql"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"snippetbox.proj.net/internal/config"
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

func NewTestDB(t *testing.T) *sql.DB {
	workDir := "/Users/alexeytarasov/Desktop/golang/src/books/lets-go/snippetbox"
	cfg, err := config.Load(workDir + "/config/local_tests.yaml")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("make", "migrate-test", "direction=up")
	cmd.Dir = workDir
	// output, _ := cmd.CombinedOutput()
	// println(string(output))
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		cmd = exec.Command("make", "migrate-test", "direction=down", "flags=-all")
		cmd.Dir = workDir
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	return db
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
		ID:       rand.Intn(1000) + 1,
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: []byte(faker.Password()),
		Created:  time.Now(),
		IsActive: true,
	}
}

func GetDummySnippet() *models.Snippet {
	return &models.Snippet{
		ID:      rand.Intn(1000) + 1,
		Title:   faker.Word(),
		Content: faker.Sentence(),
		Created: time.Now(),
		Expires: time.Now().Add(time.Hour * 24),
	}
}