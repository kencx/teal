package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func (s *Server) Run(port string) {
	s.Server.Addr = port
	s.Logger.Println("[INFO] Server listening on", port)

	go func() {
		err := s.Server.ListenAndServe()
		if err != nil {
			s.Logger.Fatalf("[ERROR]", err)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	sig := <-sigChan
	s.Logger.Println("[INFO] Received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Server.Shutdown(tc)
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
