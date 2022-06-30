package http

import (
	"github.com/kencx/teal"
)

func (s *Server) authenticate(username, password string) (bool, *teal.User, error) {
	if username == "" || password == "" {
		return false, nil, nil
	}

	user, err := s.Users.GetByUsername(username)
	if err != nil {
		return false, nil, err
	}

	// user should never be nil
	authenticated, err := user.PasswordMatches(password)
	if err != nil {
		return false, nil, err
	}
	return authenticated, user, nil
}

// func (s *Server) apiKeyValidation(key string) (bool, error) {
//
// 	user, err := s.Users.GetByAPIKey(key)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if user == nil {
// 		return false, nil
// 	}
//
// 	// check if key expired
// 	if !user.Key.Expired {
// 		return false, teal.ErrAPIKeyExpired
// 	}
//
// 	return true, nil
// }
