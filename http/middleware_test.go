package http

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kencx/teal"
	tcontext "github.com/kencx/teal/context"
	"github.com/kencx/teal/mock"
)

func TestSecureHeaders(t *testing.T) {
	next := func(rw http.ResponseWriter, r *http.Request) {
		want := map[string]string{
			"X-Frame-Options":  "deny",
			"X-XSS-Protection": "1; mode=block",
			"Set-Cookie":       "Secure; HttpOnly",
		}

		got := make(map[string]string)
		for k := range want {
			got[k] = rw.Header().Get(k)
		}
		assertObjectEqual(t, got, want)
	}

	r, err := http.NewRequest(http.MethodGet, "/api/", nil)
	checkErr(t, err)
	rw := httptest.NewRecorder()
	testServer.secureHeaders(http.HandlerFunc(next)).ServeHTTP(rw, r)
}

func TestBasicAuth(t *testing.T) {
	testUser1.SetPassword(inputTestUser1.Password)
	testServer.Users = &mock.UserStore{
		GetUserByUsernameFn: func(username string) (*teal.User, error) {
			return testUser1, nil
		},
	}

	auth := base64.StdEncoding.EncodeToString([]byte(testUser1.Username + ":" + inputTestUser1.Password))
	next := func(rw http.ResponseWriter, r *http.Request) {
		_, err := tcontext.GetUser(r.Context())
		if err == nil {
			t.Errorf("expected err: unable to get user from context")
		}
	}

	t.Run("success", func(t *testing.T) {
		next := func(rw http.ResponseWriter, r *http.Request) {
			u, err := tcontext.GetUser(r.Context())
			checkErr(t, err)

			assertEqual(t, u.Name, testUser1.Name)
			assertEqual(t, u.Username, testUser1.Username)
			assertEqual(t, u.Role, testUser1.Role)
		}
		tc := &testCase{
			url:    "/api/books/",
			method: http.MethodGet,
			data:   nil,
			params: nil,
			fn:     next,
		}

		rw, err := basicAuthTestResponse(t, tc, testServer, auth)
		checkErr(t, err)
		assertEqual(t, rw.Code, http.StatusOK)
	})

	t.Run("no auth headers", func(t *testing.T) {
		tc := &testCase{
			url:    "/api/books/",
			method: http.MethodGet,
			data:   nil,
			params: nil,
			fn:     next,
		}
		rw, err := basicAuthTestResponse(t, tc, testServer, "")
		checkErr(t, err)
		assertResponseError(t, rw, http.StatusUnauthorized, "no authentication headers")
	})

	t.Run("user does not exist", func(t *testing.T) {
		testServer.Users = &mock.UserStore{
			GetUserByUsernameFn: func(username string) (*teal.User, error) {
				return nil, teal.ErrDoesNotExist
			},
		}

		tc := &testCase{
			url:    "/api/books/",
			method: http.MethodGet,
			data:   nil,
			params: nil,
			fn:     next,
		}
		rw, err := basicAuthTestResponse(t, tc, testServer, auth)
		checkErr(t, err)
		assertResponseError(t, rw, http.StatusUnauthorized, "invalid credentials")
	})

	t.Run("auth failed", func(t *testing.T) {
		tc := &testCase{
			url:    "/api/books/",
			method: http.MethodGet,
			data:   nil,
			params: nil,
			fn:     next,
		}

		auth := base64.StdEncoding.EncodeToString([]byte(testUser1.Username + ":wrongPassword"))
		rw, err := basicAuthTestResponse(t, tc, testServer, auth)
		checkErr(t, err)
		assertResponseError(t, rw, http.StatusUnauthorized, "invalid credentials")
	})
}
