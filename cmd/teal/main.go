package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kencx/teal/author"
	"github.com/kencx/teal/book"
	"github.com/kencx/teal/http"
	"github.com/kencx/teal/storage"
)

func main() {

	a := NewApp()
	go a.Run(":9090")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan
	a.server.InfoLog.Println("Received terminate, shutting down...")

	a.Close()
	a.server.InfoLog.Println("Application gracefully stopped")
}

type App struct {
	db     *storage.Store
	server *http.Server
}

func NewApp() *App {
	return &App{
		db:     storage.NewStore("sqlite3"),
		server: http.NewServer(),
	}
}

func (a *App) Run(port string) error {

	a.db.Open("./test.db")
	a.server.InfoLog.Println("Database connection successfully established!")

	a.server.Books = book.NewService(a.db)
	a.server.Authors = author.NewService(a.db)

	if err := a.server.Run(port); err != nil {
		return err
	}
	return nil
}

func (a *App) Close() error {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return err
		}
		a.server.InfoLog.Println("Database connection closed")
	}
	if a.server != nil {
		if err := a.server.Close(); err != nil {
			return err
		}
		a.server.InfoLog.Println("Server connection closed")
	}
	return nil
}
