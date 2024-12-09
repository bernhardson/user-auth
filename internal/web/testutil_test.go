package web

import (
	"bytes"
	"context"
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

func newTestApplication() *Application {

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
	// Create a full URL for the request
	url := ts.URL + urlPath

	// Use ts.Client().Get to send the request directly
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add custom context values
	ctx := req.Context()
	var contextKey contextKey = "isAuthenticated"
	ctx = context.WithValue(ctx, contextKey, "1")
	req = req.WithContext(ctx)

	// Send the request using the test server's client
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read and process the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Return status code, headers, and trimmed body
	return resp.StatusCode, resp.Header, string(bytes.TrimSpace(body))
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
