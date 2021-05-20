package sync

import "fmt"

// Syncer defines an interface with a Sync method
type Syncer interface {
	Sync() error
}

type SyncAlreadyRunningError struct {
	name string
}

func NewSyncAlreadyRunningError(name string) error {
	return &SyncAlreadyRunningError{name: name}
}

func (e *SyncAlreadyRunningError) Error() string {
	return fmt.Sprintf("%s sync is already running", e.name)
}
