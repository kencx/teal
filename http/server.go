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

	getRouter := s.Router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", s.GetAllBooks)
	getRouter.HandleFunc("/{id:[0-9]+}", s.GetBook)

	postRouter := s.Router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", s.AddBook)
	// postRouter.HandleFunc("/{id:[0-9]+}/read", s.ReadBook)
	// postRouter.HandleFunc("/{id:[0-9]+}/reading", s.ReadingBook)
	// postRouter.HandleFunc("/{id:[0-9]+}/unread", s.UnreadBook)
	// postRouter.Use(s.MiddlewareBookValidation)

	// putRouter := s.Router.Methods(http.MethodPut).Subrouter()
	// putRouter.HandleFunc("/{id:[0-9]+}", s.UpdateBook)

	// deleteRouter := s.Router.Methods(http.MethodDelete).Subrouter()
	// deleteRouter.HandleFunc("/{id:[0-9]+}", s.DeleteBook)
}
