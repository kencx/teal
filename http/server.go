package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	idleTimeout      = 60 * time.Second
	readWriteTimeout = 3 * time.Second
	closeTimeout     = 5 * time.Second
)

type Server struct {
	Server  *http.Server
	Router  *mux.Router
	InfoLog *log.Logger
	ErrLog  *log.Logger

	Books   BookStore
	Authors AuthorStore
	Users   UserStore
}

func NewServer() *Server {
	s := &Server{
		Router:  mux.NewRouter(),
		InfoLog: log.New(os.Stdout, "INFO ", log.LstdFlags),
		ErrLog:  log.New(os.Stderr, "ERROR ", log.LstdFlags),
	}

	s.Server = &http.Server{
		Handler:      s.Router,
		ErrorLog:     s.ErrLog,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readWriteTimeout,
		WriteTimeout: readWriteTimeout,
	}

	s.RegisterRoutes()
	return s
}

func (s *Server) Run(port string) error {
	s.Server.Addr = port

	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() error {
	tc, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	return s.Server.Shutdown(tc)
}

func (s *Server) RegisterRoutes() {

	router := s.Router
	router.Use(s.recoverPanic)
	router.Use(s.logging)
	// router.Use(s.handleCORS)
	router.Use(s.secureHeaders)
	router.HandleFunc("/health", s.Healthcheck).Methods(http.MethodGet)
	router.HandleFunc("/api/users/", s.Register).Methods(http.MethodPost)
	// r.HandleFunc("/api/tokens/", s.NewToken).Methods(http.MethodPost)

	api := router.PathPrefix("/api").Subrouter()
	// api.Use(s.apiKeyAuth)
	api.Use(s.basicAuth)

	ur := api.PathPrefix("/users/").Subrouter()
	ur.HandleFunc("/{id:[0-9]+}/", s.GetUser).Methods(http.MethodGet)
	ur.HandleFunc("/{username}/", s.GetUserByUsername).Methods(http.MethodGet)

	br := api.PathPrefix("/books/").Subrouter()
	br.HandleFunc("/{id:[0-9]+}/", s.GetBook).Methods(http.MethodGet)
	br.HandleFunc("/{isbn}/", s.GetBookByISBN).Methods(http.MethodGet)
	br.HandleFunc("/", s.GetAllBooks).Methods(http.MethodGet)
	br.HandleFunc("/", s.AddBook).Methods(http.MethodPost)
	br.HandleFunc("/{id:[0-9]+}/", s.UpdateBook).Methods(http.MethodPut)
	br.HandleFunc("/{id:[0-9]+}/", s.DeleteBook).Methods(http.MethodDelete)

	ar := api.PathPrefix("/authors/").Subrouter()
	ar.HandleFunc("/{id:[0-9]+}/", s.GetAuthor).Methods(http.MethodGet)
	ar.HandleFunc("/{name}/", s.GetAuthorByName).Methods(http.MethodGet)
	ar.HandleFunc("/", s.GetAllAuthors).Methods(http.MethodGet)
	ar.HandleFunc("/", s.AddAuthor).Methods(http.MethodPost)
	ar.HandleFunc("/{id:[0-9]+}/", s.UpdateAuthor).Methods(http.MethodPut)
	ar.HandleFunc("/{id:[0-9]+}/", s.DeleteAuthor).Methods(http.MethodDelete)
}
