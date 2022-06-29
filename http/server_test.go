package http

import (
	"bytes"
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
	url    string
	method string
	data   []byte
	params map[string]string
	fn     func(http.ResponseWriter, *http.Request)
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

func checkErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
