package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal/http"
	"github.com/kencx/teal/storage"
)

type App struct {
	db     *storage.Store
	server *http.Server
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		db:     storage.NewStore(db),
		server: http.NewServer(),
	}
}

func (a *App) Run(port string) error {

	a.server.Books = a.db.Books
	a.server.Authors = a.db.Authors

	if err := a.server.Run(port); err != nil {
		return err
	}
	return nil
}

func (a *App) Close() error {
	db := a.db.GetDB()
	if db != nil {
		if err := storage.Close(db); err != nil {
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

func main() {

	db, err := storage.Open("./test.db")
	if err != nil {
		log.Fatal(err)
	}

	// init test data
	if err := storage.ExecFile(db, "../testdata/schema.sql"); err != nil {
		log.Fatal(err)
	}

	if err := storage.ExecFile(db, "../../storage/testdata.sql"); err != nil {
		log.Fatal(err)
	}

	app := NewApp(db)
	app.server.InfoLog.Println("Database connection successfully established!")

	go app.Run(":9090")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan
	app.server.InfoLog.Println("Received terminate, shutting down...")

	app.Close()
	app.server.InfoLog.Println("Application gracefully stopped")
}
