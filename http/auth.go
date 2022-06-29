package http

func (s *Server) authenticate(username, password string) (bool, error) {

	user, err := s.Users.GetByUsername(username)
	if err != nil {
		return false, err
	}

	authenticated, err := user.PasswordMatches(password)
	if err != nil {
		return false, err
	}

	return authenticated, nil
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
