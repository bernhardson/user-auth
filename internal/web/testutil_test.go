package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bernhardson/stub/internal/repo"
	"github.com/rs/zerolog"
)

func newTestApplication(t *testing.T) *Application {

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	logger := zerolog.New(io.Discard)
	return &Application{
		UserRepo:       &repo.MockUserRepo{},
		Logger:         &logger,
		SessionManager: sessionManager,
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize the test server as normal.
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response cookies will
	// now be stored and sent with subsequent requests when using this client.
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func post[T any](testing *testing.T, url string, t T, ts *testServer) (int, string) {
	bodyJson, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("Error marshalling user: %v", err)
	}
	rs, err := ts.Client().Post(ts.URL+url, "application/json", bytes.NewReader(bodyJson))
	if err != nil {
		fmt.Println(rs.Body)
		testing.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		testing.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, string(body)
}
