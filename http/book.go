package http

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kencx/teal"
	"github.com/kencx/teal/http/response"
	"github.com/kencx/teal/json"
)

type KeyBook struct{}

type BookService interface {
	Get(id int) (*teal.Book, error)
	GetByTitle(title string) (*teal.Book, error)
	GetAll() ([]*teal.Book, error)
	Create(ctx context.Context, b *teal.Book) (*teal.Book, error)
	// Update(id int, b *teal.Book) error
	Delete(id int) error
}

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)

	b, err := s.Books.Get(id)
	if err == sql.ErrNoRows {
		s.InfoLog.Printf("No book %d found", id)
		response.NoContent(rw, r)
		return
	}

	if err != nil {
		response.Error(rw, r, err)
		return
	}

	res, err := json.ToJSON(b)
	if err != nil {
		response.Error(rw, r, err)
		return
	}

	s.InfoLog.Printf("Book %d returned", id)
	response.OK(rw, r, res)
}

func (s *Server) GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	b, err := s.Books.GetAll()
	if err != nil {
		response.Error(rw, r, err)
		return
	}

	res, err := json.ToJSON(b)
	if err != nil {
		response.Error(rw, r, err)
		return
	}

	s.InfoLog.Printf("%d books returned", len(b))
	response.OK(rw, r, res)
}

func (s *Server) AddBook(rw http.ResponseWriter, r *http.Request) {
	book := r.Context().Value(KeyBook{}).(teal.Book)

	result, err := s.Books.Create(r.Context(), &book)
	if err != nil {
		s.ErrLog.Print(err)
		response.Error(rw, r, err)
		return
	}

	body, err := json.ToJSON(result)
	if err != nil {
		s.ErrLog.Print(err)
		response.Error(rw, r, err)
		return
	}
	s.InfoLog.Printf("Book %v created", result)
	response.Created(rw, r, body)
}

// func (s *Server) UpdateBook(rw http.ResponseWriter, r *http.Request) {
// 	id := HandleId(rw, r)
// 	book := r.Context().Value(KeyBook{}).(teal.Book)
//
// 	err := s.Books.Update(id, &book)
// 	if err == teal.ErrBookNotFound {
// 		http.Error(rw, "Book not found", http.StatusNotFound)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(rw, "Book not found", http.StatusInternalServerError)
// 		return
// 	}
// 	s.ErrLog.Println("Handle PUT Book", id)
// }

// func (s *Server) DeleteBook(rw http.ResponseWriter, r *http.Request) {
// 	id := HandleId(rw, r)
//
// 	err := s.Books.Delete(id)
// 	if err == teal.ErrBookNotFound {
// 		http.Error(rw, "Book not found", http.StatusNotFound)
// 		return
// 	}
// 	if err != nil {
// 		s.ErrLog.Println("Book not found")
// 		http.Error(rw, "Book not found", http.StatusInternalServerError)
// 		return
// 	}
// 	s.InfoLog.Println("Handle DELETE Book", id)
// }

func HandleId(rw http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	// TODO FIX!!
	if id == 0 {
		return 0
	}
	if err != nil {
		response.Error(rw, r, fmt.Errorf("unable to process id"))
		return -1
	}
	return id
}

func (s Server) MiddlewareBookValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		book := teal.Book{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(rw, r, err)
			return
		}
		err = json.FromJSON(body, &book)
		if err != nil {
			response.Error(rw, r, err)
			return
		}

		// TODO should this be in middleware or service layer?
		err = book.Validate()
		if err != nil {
			s.ErrLog.Printf("validating book failed: %v", err)
			response.Error(rw, r, fmt.Errorf("validation for book failed: %v", err))
			return
		}

		// add book to context
		ctx := context.WithValue(r.Context(), KeyBook{}, book)
		r = r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(rw, r)
	})
}
