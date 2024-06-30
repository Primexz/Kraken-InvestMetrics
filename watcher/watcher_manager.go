package watcher

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type WatcherManager struct {
	watchers []Watcher
	log      *logrus.Entry

	//we need a mutex to initialize the watchers (running in a goroutine) in a sequential order
	mu sync.Mutex
}

func NewWatcherManager() *WatcherManager {
	return &WatcherManager{
		log: logrus.WithFields(logrus.Fields{
			"prefix": "watcher_manager",
		}),
	}
}

func (wm *WatcherManager) RegisterWatcher(w Watcher) {
	wm.watchers = append(wm.watchers, w)
	wm.log.Debugf("Registered watcher: %T", w)
}

// StartAllWatchers starts all registered watchers
func (wm *WatcherManager) StartAllWatchers() {
	wm.log.Info("Starting all watcher clients")

	for _, w := range wm.watchers {
		go func(w Watcher) {
			for {
				wm.mu.Lock()

				start := time.Now()

				wm.log.Infof("Updating data for watcher: %T", w)
				w.UpdateData()
				wm.log.Infof("Updated data for watcher: %T in %v", w, time.Since(start))

				//we cannot use a defer here because we need to unlock the mutex before the sleep
				wm.mu.Unlock()

				<-time.After(w.GetInterval())
			}
		}(w)

		wm.log.Debugf("Started watcher: %T", w)
	}
}
