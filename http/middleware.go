package http

import (
	"net/http"

	"github.com/kencx/teal"
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

func (s *Server) basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		username, pass, ok := r.BasicAuth()
		if !ok {
			response.Unauthorized(rw, r, teal.ErrNoAuthHeader)
			return
		}

		authenticated, err := s.authenticate(username, pass)
		if err != nil {
			response.InternalServerError(rw, r, err)
			return
		}
		if !authenticated {
			response.Unauthorized(rw, r, teal.ErrInvalidCreds)
			return
		}

		// TODO save user to context
		next.ServeHTTP(rw, r)
	})
}

// func (s *Server) apiKeyAuth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//
// 		authenticated, err := s.apiKeyValidation(apiKey)
// 		// TODO check api token expired
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, teal.ErrAPIKeyExpired):
// 				response.Unauthorized(rw, r, err)
// 			default:
// 				response.InternalServerError(rw, r, err)
// 				return
// 			}
// 		}
// 		if !authenticated {
// 			response.Unauthorized(rw, r, teal.ErrInvalidCreds)
// 			return
// 		}
//
// 		// TODO save user to context
//
// 		next.ServeHTTP(rw, r)
// 	})
// }
