package storage

import "fmt"

func RunMigration(storage Provider) error {
	if err := storage.EnsureBucket(); err != nil {
		return fmt.Errorf("storage migration failed: %w", err)
	}

	return nil
}
