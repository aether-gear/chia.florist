package bootstrap

import (
	"log"
	"net/http"
)

type App struct {
	infra     *Infra
	container *Container
	router    *http.ServeMux
}

func New(cfg Config) (*App, error) {
	infra, err := NewInfra(cfg)
	if err != nil {
		return nil, err
	}

	c := NewContainer(cfg, infra)
	r := NewRouter(c)

	return &App{
		infra:     infra,
		container: c,
		router:    r,
	}, nil
}

func (a *App) Close() {
	if a == nil || a.infra == nil {
		return
	}

	a.infra.Close()
}

func (a *App) Run() error {
	log.Println("service-core running on :8000")
	return http.ListenAndServe(":8000", a.router)
}
