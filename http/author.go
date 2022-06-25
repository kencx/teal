package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/request"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type AuthorStore interface {
	Get(id int) (*teal.Author, error)
	GetAll() ([]*teal.Author, error)
	Create(ctx context.Context, b *teal.Author) (*teal.Author, error)
	Update(ctx context.Context, id int, b *teal.Author) (*teal.Author, error)
	Delete(ctx context.Context, id int) error
}

func (s *Server) GetAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	a, err := s.Authors.Get(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d not found", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(a)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d returned", id)
	response.OK(rw, r, res)
}

func (s *Server) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {

	b, err := s.Authors.GetAll()
	if err == teal.ErrNoRows {
		s.InfoLog.Println("No authors found")
		response.NoContent(rw, r)
		return

	} else if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	res, err := util.ToJSON(b)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("%d authors returned", len(b))
	response.OK(rw, r, res)
}

func (s *Server) AddAuthor(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var author teal.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	// validate payload
	verr := author.Validate()
	if verr != nil {
		// log
		response.ValidationError(rw, r, []*teal.ValidationError{verr})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := s.Authors.Create(ctx, &author)
	if err != nil {
		s.ErrLog.Print(err)
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("Author %v created", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var author teal.Author
	err := request.Read(rw, r, &author)
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	verr := author.Validate()
	if verr != nil {
		// log
		response.ValidationError(rw, r, []*teal.ValidationError{verr})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := s.Authors.Update(ctx, id, &author)
	if err == teal.ErrDoesNotExist {
		response.InternalServerError(rw, r, err)
		return
	}
	if err != nil {
		response.InternalServerError(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %v updated", result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteAuthor(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	err := s.Authors.Delete(r.Context(), id)
	if err == teal.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Author %d deleted", id)
	response.OK(rw, r, nil)
}
