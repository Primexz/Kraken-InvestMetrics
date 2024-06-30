package watcher

import (
	"sync"
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/util"
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
		log: util.LoggerWithPrefix("watcher_manager"),
	}
}

func (wm *WatcherManager) RegisterWatcher(w Watcher) {
	wm.watchers = append(wm.watchers, w)
	wm.log.Debugf("Registered watcher: %T", w)
}

// StartAllWatchers starts all registered watchers
func (wm *WatcherManager) StartAllWatchers() {
	wm.log.Info("Starting all watcher clients")

	var (
		wg       sync.WaitGroup
		finished = make(chan struct{})
	)

	for _, w := range wm.watchers {
		wg.Add(1)

		go func() {
			for {
				wm.mu.Lock()

				start := time.Now()

				wm.log.Infof("Updating data for watcher: %T", w)
				w.UpdateData()
				wm.log.Infof("Updated data for watcher: %T in %v", w, time.Since(start))

				wm.mu.Unlock()
				wg.Done()

				<-time.After(w.GetInterval())
			}
		}()

		wm.log.Debugf("Started watcher: %T", w)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	<-finished
	wm.log.Info("All watchers have completed their initial updates")
}
