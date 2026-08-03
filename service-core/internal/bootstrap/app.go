package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	dependencies *Dependency
	container    *Container
	router       *chi.Mux
	cfg          Config
	cancelSync   context.CancelFunc
}

func New(cfg Config) (*App, error) {
	dependencies, err := NewDependency(cfg)
	if err != nil {
		return nil, err
	}

	c := NewContainer(cfg, dependencies)

	// Perform startup payment method synchronization
	if err := c.SyncPaymentMethods(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to sync payment methods: %w", err)
	}

	r := NewRouter(c)

	return &App{
		dependencies: dependencies,
		container:    c,
		router:       r,
		cfg:          cfg,
	}, nil
}

func (a *App) Close() {
	// Stop the reconciliation job first so it doesn't race
	// against the DB connection being closed.
	if a.cancelSync != nil {
		a.cancelSync()
	}

	if a == nil || a.dependencies == nil {
		return
	}

	a.dependencies.Close()
}

func (a *App) Run() error {
	// Start the payment reconciliation job in the background.
	// It shuts down cleanly when the context is cancelled (in Close).
	if a.dependencies.PaymentSyncJob != nil {
		syncCtx, cancel := context.WithCancel(context.Background())
		a.cancelSync = cancel

		go a.dependencies.PaymentSyncJob.Start(syncCtx)
	}

	addr := a.cfg.App.Host + ":" + a.cfg.App.Port
	log.Printf("service-core running on %s\n", addr)
	return http.ListenAndServe(addr, a.router)
}
