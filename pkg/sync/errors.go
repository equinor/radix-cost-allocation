package sync

import (
	"errors"
	"fmt"
)

// Syncer defines an interface with a Sync method
type Syncer interface {
	Sync() error
}

var ErrSyncAlreadyRunning = errors.New("sync is already running")

func NewSyncAlreadyRunningError(name string) error {
	return fmt.Errorf("%w: %s", ErrSyncAlreadyRunning, name)
}
