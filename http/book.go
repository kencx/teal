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

func (s *Server) GetBook(rw http.ResponseWriter, r *http.Request) {
	id := HandleId(rw, r)
	if id == -1 {
		return
	}

	b, err := s.Store.RetrieveBookWithID(id)
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
	b, err := s.Store.RetrieveAllBooks()

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

	result, err := s.Store.CreateBook(ctx, &book)
	if err != nil {
		s.ErrLog.Print(err)
		response.Error(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.Error(rw, r, err)
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
		response.Error(rw, r, err)
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

	result, err := s.Store.UpdateBook(ctx, id, &book)
	if err == teal.ErrDoesNotExist {
		response.Error(rw, r, err)
		return
	}
	if err != nil {
		response.Error(rw, r, err)
		return
	}

	body, err := util.ToJSON(result)
	if err != nil {
		s.ErrLog.Println(err)
		response.Error(rw, r, err)
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

	err := s.Store.DeleteBook(r.Context(), id)
	if err == teal.ErrDoesNotExist {
		http.Error(rw, "Book not found", http.StatusNotFound)
		return
	}

	if err != nil {
		s.ErrLog.Println(err)
		response.Error(rw, r, err)
		return
	}

	s.InfoLog.Printf("Book %d deleted", id)
	response.OK(rw, r, nil)
}

func HandleId(rw http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(rw, r, fmt.Errorf("unable to process id: %v", err))
		return -1
	}

	return id
}
