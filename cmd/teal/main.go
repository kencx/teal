package main

import (
	"github.com/kencx/teal/pkg/http"
	"github.com/kencx/teal/pkg/storage"
)

func main() {

	a := NewApp()
	a.Run(":9090")
}

type App struct {
	DB         *storage.DB
	HTTPServer *http.Server
}

func NewApp() *App {
	return &App{
		DB:         storage.NewDB("sqlite3"),
		HTTPServer: http.NewServer(),
	}
}

func (a *App) Run(port string) {

	a.DB.Open("./test.db")
	a.HTTPServer.Logger.Println("[INFO] Database connection successfully established!")

	a.HTTPServer.BS = a.DB
	a.HTTPServer.AS = a.DB

	a.HTTPServer.Run(port)
}
