package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/kencx/teal/http"
	"github.com/kencx/teal/storage"
)

type config struct {
	port int
	env  string
	dsn  string
}

type App struct {
	config config
	db     *storage.Store
	server *http.Server
}

func NewApp(config config, db *sqlx.DB) *App {
	return &App{
		config: config,
		db:     storage.NewStore(db),
		server: http.NewServer(),
	}
}

func (a *App) Run() error {
	a.server.Books = a.db.Books
	a.server.Authors = a.db.Authors

	a.server.InfoLog.Printf("Starting %s server on :%d", a.config.env, a.config.port)
	if err := a.server.Run(fmt.Sprintf(":%d", a.config.port)); err != nil {
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

	var config config

	flag.IntVar(&config.port, "port", 9090, "API Server Port")
	flag.StringVar(&config.env, "env", "dev", "Environment (dev|staging|prod)")

	// flag.StringVar(&config.dsn, "dsn", "postgres://teal:password1@localhost/teal", "PostgreSQL DSN")
	// flag.StringVar(&config.dsn, "dsn", os.Getenv("TEAL_POSTGRES_DSN"), "PostgreSQL DSN")

	flag.Parse()

	db, err := storage.Open("./test.db")
	// db, err := storage.Open(config.dsn)
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

	app := NewApp(config, db)
	app.server.InfoLog.Println("Database connection successfully established!")

	go app.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	app.server.InfoLog.Printf("Received signal %s, shutting down...", s.String())

	app.Close()
	app.server.InfoLog.Println("Application gracefully stopped")
}
