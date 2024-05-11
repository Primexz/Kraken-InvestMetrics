package watcher

import (
	watcherClient "github.com/Primexz/Kraken-InvestMetrics/updateWatchers/watchers"
	"github.com/Primexz/Kraken-InvestMetrics/xPub"
	"github.com/primexz/KrakenDCA/logger"
)

var (
	XPubWatcher   *watcherClient.XPubWatcher
	KrakenWatcher *watcherClient.KrakenWatcher

	log *logger.Logger
)

func init() {
	log = logger.NewLogger("watcher")
}

func BootstrapWatchers() {
	log.Info("Boostrapping watchers..")

	if xPub.IsXPub() {
		XPubWatcher = watcherClient.NewXPubWatcher()
		XPubWatcher.StartRoutine()
	}

	KrakenWatcher = watcherClient.NewKrakenWatcher()
	KrakenWatcher.StartRoutine()
}
