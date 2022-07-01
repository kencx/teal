package http

import (
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/request"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
	"github.com/kencx/teal/validator"
)

type AuthorStore interface {
	Get(id int64) (*teal.Author, error)
	GetByName(name string) (*teal.Author, error)
	GetAll() ([]*teal.Author, error)
	Create(b *teal.Author) (*teal.Author, error)
	Update(id int64, b *teal.Author) (*teal.Author, error)
	Delete(id int64) error
}

func (s *Server) GetAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	a, err := s.Authors.Get(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d retrieved: %v", id, a)
	response.OK(rw, r, res)
}

func (s *Server) GetAuthorByName(rw http.ResponseWriter, r *http.Request) {
	name := HandleString("name", r)

	a, err := s.Authors.GetByName(name)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %q does not exist", name)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %q retrieved: %v", name, a)
	response.OK(rw, r, res)
}

func (s *Server) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {

	a, err := s.Authors.GetAll()
	if err == teal.ErrNoRows {
		s.InfoLog.Println("No authors retrieved")
		response.NoContent(rw, r)
		return

	} else if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(response.Envelope{"authors": a})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("%d authors retrieved: %v", len(a), a)
	response.OK(rw, r, res)
}

func (s *Server) AddAuthor(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var author teal.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	// validate payload
	v := validator.New()
	author.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.Authors.Create(&author)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("New author created: %v", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var author teal.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	v := validator.New()
	author.Validate(v)
	if !v.Valid() {
		response.ValidationError(rw, r, v.Errors)
		return
	}

	result, err := s.Authors.Update(id, &author)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(response.Envelope{"authors": result})
	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d updated: %v", id, result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleInt64("id", rw, r)
	if id == -1 {
		return
	}

	err := s.Authors.Delete(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d does not exist", id)
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrLog.Printf("err: %v", err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d deleted", id)
	response.OK(rw, r, nil)
}
