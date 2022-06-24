package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {

	var user *teal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.HashedPassword), 12)
	if err != nil {
		response.Error(w, r, err)
	}

	// validate payload
	verrs := user.Validate()
	if len(verrs) > 0 {
		// log
		response.ValidationError(w, r, verrs)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := s.Users.Create(ctx, user)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.Error(w, r, err)
		return
	}

	s.InfoLog.Printf("User %v created", result)
	response.Created(w, r, body)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) checkUser(username, password string) bool {

	user, _ := s.Users.GetByUsername(username)
	if string(user.HashedPassword) == password {
		return true
	}
	return false
}
