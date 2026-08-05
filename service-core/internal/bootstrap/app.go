package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	cfg          Config
	dependencies *Dependency
	container    *Container
	scheduler    *Scheduler
	router       *chi.Mux
	cancelJobs   context.CancelFunc
}

func New(cfg Config) (*App, error) {
	dependencies, err := NewDependency(cfg)
	if err != nil {
		return nil, err
	}

	c := NewContainer(cfg, dependencies)
	r := NewRouter(c)
	scheduler := NewScheduler(cfg, c, c.Logger)

	initializer := NewInitializer(c)
	if err := initializer.Run(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to run startup initializer: %w", err)
	}

	return &App{
		cfg:          cfg,
		dependencies: dependencies,
		container:    c,
		scheduler:    scheduler,
		router:       r,
	}, nil
}

func (a *App) Close() {
	// Stop background jobs first so they don't race
	// against the DB connection being closed.
	if a.cancelJobs != nil {
		a.cancelJobs()
	}

	if a == nil || a.dependencies == nil {
		return
	}

	a.dependencies.Close()
}

func (a *App) Run() error {
	// Start background jobs.
	//
	// The background jobs will shut down cleanly when the context is cancelled (in Close).
	jobCtx, cancel := context.WithCancel(context.Background())
	a.cancelJobs = cancel

	if a.scheduler != nil {
		a.scheduler.Start(jobCtx)
	}

	addr := a.cfg.App.Host + ":" + a.cfg.App.Port
	log.Printf("service-core running on %s\n", addr)
	return http.ListenAndServe(addr, a.router)
}
