package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/util"
)

type BookStore interface {
	Get(id int) (*teal.Book, error)
	GetByTitle(title string) (*teal.Book, error)
	GetAll() ([]*teal.Book, error)
	Create(ctx context.Context, b *teal.Book) (*teal.Book, error)
	Update(ctx context.Context, id int, b *teal.Book) (*teal.Book, error)
	Delete(ctx context.Context, id int) error

	GetByAuthor(name string) ([]*teal.Book, error)
}

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	b, err := s.Books.Get(id)
	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Book %d not found", id)
		response.NotFound(rw, r, err)
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

	s.InfoLog.Printf("Book %d returned", id)
	response.OK(rw, r, res)
}

func hasQueryParam(param string, r *http.Request) bool {
	p := r.URL.Query().Get(param)
	return p != ""
}

func (s *Server) GetAllBooks(rw http.ResponseWriter, r *http.Request) {

	var b []*teal.Book
	var err error

	if hasQueryParam("author", r) {
		b, err = s.Books.GetByAuthor(r.URL.Query().Get("author"))
	} else {
		b, err = s.Books.GetAll()
	}

	if err == teal.ErrNoRows {
		s.InfoLog.Println("No books found")
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

	s.InfoLog.Printf("%d books returned", len(b))
	response.OK(rw, r, res)
}

func (s *Server) AddBook(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var book teal.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	// validate payload
	verrs := book.Validate()
	if len(verrs) > 0 {
		// log
		response.ValidationError(rw, r, verrs)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := s.Books.Create(ctx, &book)
	if err != nil {
		s.ErrLog.Print(err)
		response.BadRequest(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}
	s.InfoLog.Printf("Book %v created", result)
	response.Created(rw, r, body)
}

func (s *Server) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	// marshal payload to struct
	var book teal.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		response.BadRequest(rw, r, err)
		return
	}

	// validate payload
	// PUT should require all fields
	verrs := book.Validate()
	if len(verrs) > 0 {
		// log
		response.ValidationError(rw, r, verrs)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := s.Books.Update(ctx, id, &book)
	if err == teal.ErrDoesNotExist {
		response.NotFound(rw, r, err)
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

	s.InfoLog.Printf("Book %v updated", result)
	response.OK(rw, r, body)
}

func (s *Server) DeleteBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	err := s.Books.Delete(r.Context(), id)
	if err == teal.ErrDoesNotExist {
		response.NotFound(rw, r, err)
		return
	}

	if err != nil {
		s.ErrLog.Println(err)
		response.InternalServerError(rw, r, err)
		return
	}

	s.InfoLog.Printf("Book %d deleted", id)
	response.OK(rw, r, nil)
}

func HandleId(rw http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.BadRequest(rw, r, fmt.Errorf("unable to process id: %v", err))
		return -1
	}

	return id
}
