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
	idleTimeout      = 120 * time.Second
	readWriteTimeout = 1 * time.Second
	closeTimeout     = 5 * time.Second
)

type Server struct {
	Server  *http.Server
	Router  *mux.Router
	InfoLog *log.Logger
	ErrLog  *log.Logger

	Books   BookService
	Authors AuthorService
}

func NewServer() *Server {
	s := &Server{
		Router:  mux.NewRouter(),
		InfoLog: log.New(os.Stdout, "INFO ", log.LstdFlags),
		ErrLog:  log.New(os.Stdout, "ERROR ", log.LstdFlags),
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

	s.InfoLog.Println("Server listening on", port)
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

	s.Router.HandleFunc("/health", s.Healthcheck).Methods(http.MethodGet)
	api := s.Router.PathPrefix("/api").Subrouter()

	br := api.PathPrefix("/books").Subrouter()
	ar := api.PathPrefix("/authors").Subrouter()

	br.HandleFunc("/{id:[0-9]+}", s.GetBook).Methods(http.MethodGet)
	br.HandleFunc("/", s.GetAllBooks).Methods(http.MethodGet)
	br.HandleFunc("/{id:[0-9]+}/", s.UpdateBook).Methods(http.MethodPut)
	br.HandleFunc("/{id:[0-9]+}/", s.DeleteBook).Methods(http.MethodDelete)
	br.HandleFunc("/", s.AddBook).Methods(http.MethodPost)

	ar.HandleFunc("/{id:[0-9]+}/", s.GetAuthor).Methods(http.MethodGet)
	ar.HandleFunc("/", s.GetAllAuthors).Methods(http.MethodGet)
	ar.HandleFunc("/{id:[0-9]+}/", s.UpdateAuthor).Methods(http.MethodPut)
	ar.HandleFunc("/{id:[0-9]+}/", s.DeleteAuthor).Methods(http.MethodDelete)
	ar.HandleFunc("/", s.AddAuthor).Methods(http.MethodPost)
}
