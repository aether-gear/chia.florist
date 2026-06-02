package bootstrap

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	dependencies *Dependency
	container    *Container
	router       *chi.Mux
}

func New(cfg Config) (*App, error) {
	dependencies, err := NewDependency(cfg)
	if err != nil {
		return nil, err
	}

	c := NewContainer(cfg, dependencies)
	r := NewRouter(c)

	return &App{
		dependencies: dependencies,
		container:    c,
		router:       r,
	}, nil
}

func (a *App) Close() {
	if a == nil || a.dependencies == nil {
		return
	}

	a.dependencies.Close()
}

func (a *App) Run() error {
	log.Println("service-core running on :8000")
	return http.ListenAndServe(":8000", a.router)
}
