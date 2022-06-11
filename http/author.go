package http

import (
	"context"
	"net/http"

	"github.com/kencx/teal"
	"github.com/kencx/teal/json"
)

type AuthorService interface {
	Get(id int) (*teal.Author, error)
	GetAll() ([]*teal.Author, error)
	Create(ctx context.Context, b *teal.Author) (int, error)
	// Update(id int, b *teal.Author) error
	Delete(ctx context.Context, id int) error
}

func (s *Server) GetAllAuthors(rw http.ResponseWriter, r *http.Request) {
	b, err := s.Authors.GetAll()
	if err != nil {
		s.ErrLog.Print(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	res, err := json.ToJSON(b)
	if err != nil {
		s.ErrLog.Print(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(res) == 0 {
		rw.Write([]byte("No books added"))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}
