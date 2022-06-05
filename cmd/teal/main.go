package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kencx/teal/pkg/http"
	"github.com/kencx/teal/pkg/storage"
)

func main() {

	a := NewApp()
	go a.Run(":9090")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan
	a.HTTPServer.Logger.Println("[INFO] Received terminate, shutting down...")

	a.Close()
	a.HTTPServer.Logger.Println("[INFO] Application gracefully stopped")
}

type App struct {
	DB         *storage.Repository
	HTTPServer *http.Server
}

func NewApp() *App {
	return &App{
		DB:         storage.NewRepository("sqlite3"),
		HTTPServer: http.NewServer(),
	}
}

func (a *App) Run(port string) error {

	a.DB.Open("./test.db")
	a.HTTPServer.Logger.Println("[INFO] Database connection successfully established!")

	a.HTTPServer.Books = a.DB
	// a.HTTPServer.Authors = a.DB

	if err := a.HTTPServer.Run(port); err != nil {
		return err
	}
	return nil
}

func (a *App) Close() error {
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			return err
		}
		a.HTTPServer.Logger.Println("[INFO] Database connection closed")
	}
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Close(); err != nil {
			return err
		}
		a.HTTPServer.Logger.Println("[INFO] Server connection closed")
	}
	return nil
}
