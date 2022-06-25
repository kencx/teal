package http

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	s := Server{
		InfoLog: testInfoLog,
		ErrLog:  testErrLog,
	}

	w, err := getResponse("/health", s.Healthcheck)
	checkErr(t, err)

	var env map[string]health
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["health"]
	assertEqual(t, got.Version, "v1.0")
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}
