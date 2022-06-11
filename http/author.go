package http

import (
	"context"
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type AuthorService interface {
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

	b, err := s.Authors.Get(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Author %d not found", id)
		response.NotFound(rw, r, err)
		return

	} else if err != nil {
		response.Error(rw, r, err)
		return
	}

	res, err := util.ToJSON(b)
	if err != nil {
		response.Error(rw, r, err)
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
		response.Error(rw, r, err)
		return
	}

	res, err := util.ToJSON(b)
	if err != nil {
		response.Error(rw, r, err)
		return
	}

	s.InfoLog.Printf("%d authors returned", len(b))
	response.OK(rw, r, res)
}
