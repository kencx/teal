package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kencx/teal/pkg"
)

type KeyBook struct{}

func (s *Server) GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	b, err := s.BS.GetAllBooks()
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	res, err := ToJSON(b)
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
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

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)

	b, err := s.BS.GetBook(id)
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := ToJSON(b)
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}

func (s *Server) AddBook(rw http.ResponseWriter, r *http.Request) {
	book := r.Context().Value(KeyBook{}).(pkg.Book)
	id, err := s.BS.CreateBook(&book)
	if err != nil {
		s.Logger.Printf("[ERROR] %v", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "Book %d created", id)
}

func (s *Server) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	book := r.Context().Value(KeyBook{}).(pkg.Book)

	err := s.BS.UpdateBook(id, &book)
	if err == pkg.ErrBookNotFound {
		http.Error(rw, "Book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Book not found", http.StatusInternalServerError)
		return
	}
	s.Logger.Println("Handle PUT Book", id)
}

func (s *Server) DeleteBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)

	err := s.BS.DeleteBook(id)
	if err == pkg.ErrBookNotFound {
		http.Error(rw, "Book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Book not found", http.StatusInternalServerError)
		return
	}
	s.Logger.Println("Handle DELETE Book", id)
}

func HandleId(rw http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return -1
	}
	return id
}

func (s Server) MiddlewareBookValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		book := pkg.Book{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = FromJSON(body, &book)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = book.Validate()
		if err != nil {
			s.Logger.Println("[ERROR] validating book", err)
			http.Error(rw, fmt.Sprintf("Error validating book: %s", err), http.StatusBadRequest)
			return
		}

		// add book to context
		ctx := context.WithValue(r.Context(), KeyBook{}, book)
		r = r.WithContext(ctx)

		// call next handler
		next.ServeHTTP(rw, r)
	})
}
