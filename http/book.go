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

type BookService interface {
	Get(id int) (*teal.Book, error)
	GetByTitle(title string) (*teal.Book, error)
	GetAll() ([]*teal.Book, error)
	Create(ctx context.Context, b *teal.Book) (*teal.Book, error)
	// Update(id int, b *teal.Book) error
	Delete(ctx context.Context, id int) error
}

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	b, err := s.Books.Get(id)

	if err == teal.ErrDoesNotExist {
		s.InfoLog.Printf("Book %d not found", id)
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

	s.InfoLog.Printf("Book %d returned", id)
	response.OK(rw, r, res)
}

func (s *Server) GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	b, err := s.Books.GetAll()

	if err == teal.ErrNoRows {
		s.InfoLog.Println("No books found")
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

	s.InfoLog.Printf("%d books returned", len(b))
	response.OK(rw, r, res)
}

func (s *Server) AddBook(rw http.ResponseWriter, r *http.Request) {

	// marshal payload to struct
	var book teal.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		response.Error(rw, r, err)
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
		response.Error(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
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
