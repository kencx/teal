package http

import (
	"errors"
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/request"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
	"github.com/kencx/teal/validator"
)

type UserStore interface {
	Get(id int64) (*teal.User, error)
	GetByUsername(username string) (*teal.User, error)
	Create(u *teal.User) (*teal.User, error)
	Update(id int64, b *teal.User) (*teal.User, error)
	Delete(id int64) error
}

func (s *Server) Register(rw http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string
		Username string
		Email    string
		Password string
	}

	err := request.Read(rw, r, &input)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	user := teal.User{
		Name:     input.Name,
		Username: input.Username,
	}

	err = user.HashedPassword.Set(input.Password)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	v := validator.New()
	user.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.Users.Create(&user)
	if err != nil {
		switch {
		case errors.Is(err, teal.ErrDuplicateUsername):
			v.AddError("username", "this username already exists")
			response.ValidationError(rw, r, v.Errors)
			return
		default:
			response.InternalServerError(rw, r, err)
			return
		}
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

func (s *Server) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var user teal.User
	err := request.Read(rw, r, &user)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	// TODO handle change password

	// validate payload
	// PUT should require all fields
	v := validator.New()
	user.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.Users.Update(id, &user)
	if err == teal.ErrDoesNotExist {
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"users": result})
	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %v updated", result)
	response.OK(rw, r, body)

}

func (s *Server) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	err := s.Users.Delete(id)
	if err == teal.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %d deleted", id)
	response.OK(rw, r, nil)

}
