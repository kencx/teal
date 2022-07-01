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

func (s *Server) GetUser(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	u, err := s.Users.Get(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("User %d does not exist", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"users": u})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %d retrieved: %v", id, u)
	response.OK(rw, r, res)
}

func (s *Server) GetUserByUsername(rw http.ResponseWriter, r *http.Request) {
	username := HandleString("username", r)

	u, err := s.Users.GetByUsername(username)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("User %q does not exist", username)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"users": u})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %q retrieved: %v", username, u)
	response.OK(rw, r, res)
}

func (s *Server) Register(rw http.ResponseWriter, r *http.Request) {

	var input teal.InputUser
	err := request.Read(rw, r, &input)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.BadRequest(rw, r, err)
		return
	}

	v := validator.New()
	input.Validate(v)

	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	user := teal.User{
		Name:     input.Name,
		Username: input.Username,
	}

	err = user.SetPassword(input.Password)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
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
			s.ErrLog.Printf("err: %v", err)
			response.InternalServerError(rw, r, err)
			return
		}
	}

	body, err := util.ToJSON(response.Envelope{"users": result})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("New user registered: %v", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	var input teal.InputUser
	err := request.Read(rw, r, &input)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	v := validator.New()
	input.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	user := teal.User{
		Name:     input.Name,
		Username: input.Username,
	}

	// TODO handle change password

	result, err := s.Users.Update(id, &user)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("User %d does not exist", id)
		response.NoContent(rw, r)
		return
	}
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"users": result})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %d updated: %v", id, result)
	response.OK(rw, r, body)

}

func (s *Server) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.Users.Delete(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("User %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("User %d deleted", id)
	response.OK(rw, r, nil)
}
