package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kencx/teal/pkg"
)

var (
	idleTimeout      = 120 * time.Second
	readWriteTimeout = 1 * time.Second
)

type Server struct {
	Server *http.Server
	Router *mux.Router
	Logger *log.Logger

	BS pkg.BookService
	AS pkg.AuthorService
}

func NewServer() *Server {
	s := &Server{
		Server: &http.Server{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
		},
		Router: mux.NewRouter(),
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	s.Server.Handler = s.Router
	s.RegisterRoutes()

	return s
}

func (s *Server) Run(port string) error {
	s.Server.Addr = port

	err := s.Server.ListenAndServe()
	if err != nil {
		return err
	}

	s.Logger.Println("[INFO] Server listening on", port)
	return nil
}

func (s *Server) Close() error {
	tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	postRouter.Use(s.MiddlewareBookValidation)

	putRouter := s.Router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", s.UpdateBook)
	putRouter.Use(s.MiddlewareBookValidation)

	deleteRouter := s.Router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", s.DeleteBook)
	putRouter.Use(s.MiddlewareBookValidation)
}
