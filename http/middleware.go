package http

import (
	"errors"
	"net/http"

	"github.com/kencx/teal/http/response"
)

func (s *Server) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// cookies
		w.Header().Set("Set-Cookie", "Secure; HttpOnly")

		next.ServeHTTP(w, r)
	})
}

// func (s *Server) handleCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
// 	})
// }

var (
	errNoAuthHeader = errors.New("no authentication headers")
	errInvalidCreds = errors.New("invalid username or password")
)

func (s *Server) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()
		if !ok {
			response.Unauthorized(w, r, errNoAuthHeader)
			return
		}

		if !s.checkUser(user, pass) {
			response.Unauthorized(w, r, errInvalidCreds)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func (s *Server) apiKeyAuth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
// 		next.ServeHTTP(w, r)
// 	})
// }
