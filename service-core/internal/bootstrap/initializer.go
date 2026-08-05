package bootstrap

import (
	"context"
	"fmt"
)

// Initializer handles pre-flight application startup tasks
// before traffic is served.
type Initializer struct {
	c *Container
}

func NewInitializer(c *Container) *Initializer {
	return &Initializer{c: c}
}

func (i *Initializer) Run(ctx context.Context) error {
	if i == nil || i.c == nil {
		return nil
	}

	if err := i.c.SyncPaymentMethods.Execute(ctx); err != nil {
		return fmt.Errorf("failed to sync payment methods on startup: %w", err)
	}

	return nil
}
