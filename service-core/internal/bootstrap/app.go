package bootstrap

import (
	"log"
	"net/http"
)

type App struct {
	container *Container
	router    *http.ServeMux
}

func NewApp() *App {
	c := NewContainer()
	r := NewRouter(c)

	return &App{
		container: c,
		router:    r,
	}
}

func (a *App) Run() {
	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", a.router))
}
