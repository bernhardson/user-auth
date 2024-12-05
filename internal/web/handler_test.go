package web

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/bernhardson/stub/internal/assert"
	"github.com/bernhardson/stub/internal/models"
)

func TestUserGet(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name          string
		urlPath       string
		expectedCode  int
		expectedEmail string
	}{
		{
			name:          "valid",
			urlPath:       "/user/view/1",
			expectedCode:  http.StatusOK,
			expectedEmail: "john.doe@gmail.com",
		},
		{
			name:         "not-existent-1",
			urlPath:      "/user/view/2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "not-existent-2",
			urlPath:      "/user/view/-1",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "not-existent-3",
			urlPath:      "/user/view/1.23",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedCode)
			assert.StringContains(t, body, tt.expectedEmail)
		})
	}
}

func TestUserSignUpPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name               string
		user               models.User
		expectedHttpStatus int
		expectedMessage    string
	}{
		{
			name: "valid",
			user: models.User{
				Username: "peterson",
				Email:    "peterson@abc.de",
				Password: "12345678",
			},
			expectedHttpStatus: http.StatusOK,
			expectedMessage:    "User created successfully",
		},
		{
			name: "invalid-password",
			user: models.User{
				Username: "peterson",
				Email:    "peterson@abc.de",
				Password: "1234",
			},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "Password must have at least 8 characters.",
		},
		{
			name: "invalid-email",
			user: models.User{
				Username: "peterson",
				Email:    "petersonc.de",
				Password: "12345678",
			},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "Entered email adress is not valid.",
		},
		{
			name: "invalid-name",
			user: models.User{
				Username: "ab",
				Email:    "petersonc.de",
				Password: "12345678",
			},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "User name must have at least three characters.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualHttpStatus, body := post(t, "/user/signup", tt.user, ts)
			assert.Equal(t, actualHttpStatus, tt.expectedHttpStatus)
			assert.StringContains(t, string(body), tt.expectedMessage)

		})
	}

}

func TestUserLoginPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name                string
		body                userLoginPost
		pass                string
		expectedHttpStatus  int
		expectedMessage     string
		expectSessionCookie bool
	}{
		{
			name:                "valid",
			body:                userLoginPost{Email: "john.doe@gmail.com", Password: "jd12345678"},
			expectedHttpStatus:  http.StatusOK,
			expectedMessage:     "User logged in",
			expectSessionCookie: true,
		},
		{
			name:                "wrong-pass",
			body:                userLoginPost{Email: "john.doe@gmail.com", Password: "123"},
			expectedHttpStatus:  http.StatusUnauthorized,
			expectedMessage:     "Email or password is incorrect",
			expectSessionCookie: true,
		},
		{
			name:               "blank-email",
			body:               userLoginPost{Email: "", Password: "jd12345678"},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "This field cannot be blank",
		},
		{
			name:               "invalid-email",
			body:               userLoginPost{Email: "12312.de", Password: "jd12345678"},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "Entered email adress is not valid",
		},
		{
			name:               "blank-password",
			body:               userLoginPost{Email: "john.doe@gmail.com", Password: ""},
			expectedHttpStatus: http.StatusUnprocessableEntity,
			expectedMessage:    "This field cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualHttpStatus, body := post(t, "/user/login", tt.body, ts)
			assert.Equal(t, actualHttpStatus, tt.expectedHttpStatus)
			assert.StringContains(t, body, tt.expectedMessage)
		})
	}
}

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
