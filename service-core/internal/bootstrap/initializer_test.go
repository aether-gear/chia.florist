package bootstrap

import (
	"context"
	"testing"
)

func TestInitializer_NilContainer(t *testing.T) {
	init := NewInitializer(nil)
	if err := init.Run(context.Background()); err != nil {
		t.Errorf("expected no error for nil container, got %v", err)
	}
}
