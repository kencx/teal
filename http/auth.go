package http

import (
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/request"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
	"github.com/kencx/teal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Register(rw http.ResponseWriter, r *http.Request) {

	var user *teal.User
	err := request.Read(rw, r, &user)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	user.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.HashedPassword), 12)
	if err != nil {
		response.InternalServerError(rw, r, err)
	}

	// validate payload
	v := validator.New()
	user.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.Users.Create(user)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"user": result})
	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %v created", result)
	response.Created(rw, r, body)
}

func (s *Server) Login(rw http.ResponseWriter, r *http.Request) {

}

func (s *Server) Logout(rw http.ResponseWriter, r *http.Request) {

}

func (s *Server) checkUser(username, password string) bool {

	user, _ := s.Users.GetByUsername(username)
	if string(user.HashedPassword) == password {
		return true
	}
	return false
}
