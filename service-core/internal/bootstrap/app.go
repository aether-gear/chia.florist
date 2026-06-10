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
	cfg          Config
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
		cfg:          cfg,
	}, nil
}

func (a *App) Close() {
	if a == nil || a.dependencies == nil {
		return
	}

	a.dependencies.Close()
}

func (a *App) Run() error {
	addr := a.cfg.App.Host + ":" + a.cfg.App.Port
	log.Printf("service-core running on %s\n", addr)
	return http.ListenAndServe(addr, a.router)
}
