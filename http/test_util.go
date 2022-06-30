package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

var testServer = Server{
	InfoLog: log.New(io.Discard, "", log.LstdFlags),
	ErrLog:  log.New(io.Discard, "", log.LstdFlags),
}

type testCase struct {
	url     string
	method  string
	headers map[string]string
	data    []byte
	params  map[string]string
	fn      func(http.ResponseWriter, *http.Request)
}

func testResponse(t *testing.T, tc *testCase) (*httptest.ResponseRecorder, error) {
	t.Helper()

	req, err := http.NewRequest(tc.method, tc.url, bytes.NewReader(tc.data))
	if err != nil {
		return nil, err
	}

	rw := httptest.NewRecorder()
	if tc.params != nil {
		req = mux.SetURLVars(req, tc.params)
	}

	http.HandlerFunc(tc.fn).ServeHTTP(rw, req)
	return rw, nil
}

func middlewareTestResponse(t *testing.T, tc *testCase, fn func(next http.Handler) http.Handler) (*httptest.ResponseRecorder, error) {
	t.Helper()

	req, err := http.NewRequest(tc.method, tc.url, bytes.NewReader(tc.data))
	if err != nil {
		return nil, err
	}
	if tc.headers != nil {
		for k, v := range tc.headers {
			req.Header.Add(k, v)
		}
	}

	rw := httptest.NewRecorder()
	if tc.params != nil {
		req = mux.SetURLVars(req, tc.params)
	}

	fn(http.HandlerFunc(tc.fn)).ServeHTTP(rw, req)
	return rw, nil
}

func basicAuthTestResponse(t *testing.T, tc *testCase, auth string) (*httptest.ResponseRecorder, error) {
	tc.headers = map[string]string{"Authorization": "Basic " + auth}
	rw, err := middlewareTestResponse(t, tc, testServer.basicAuth)
	checkErr(t, err)
	return rw, nil
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertObjectEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertResponseError(t *testing.T, w *httptest.ResponseRecorder, status int, message string) {
	var env map[string]string
	err := json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, status)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
	assertEqual(t, got, message)
}

func assertValidationError(t *testing.T, w *httptest.ResponseRecorder, key, message string) {
	t.Helper()

	var env map[string]map[string]string
	err := json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["error"]
	assertEqual(t, w.Code, http.StatusUnprocessableEntity)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")

	val, ok := got[key]
	if !ok {
		t.Errorf("validation error field %q not present", key)
	}
	assertEqual(t, val, message)
}

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
