package watcher

import "time"

// Watcher is the base interface that all watchers must implement
type Watcher interface {
	UpdateData()

	GetInterval() time.Duration
}
