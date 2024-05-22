package watcher

import (
	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	watcherClient "github.com/Primexz/Kraken-InvestMetrics/updateWatchers/watchers"
	"github.com/primexz/KrakenDCA/logger"
)

var (
	XPubWatcher     *watcherClient.XPubWatcher
	KrakenWatcher   *watcherClient.KrakenWatcher
	PurchaseWatcher *watcherClient.PurchaseWatcher
	DCAWatcher      *watcherClient.DCAWatcher

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

	PurchaseWatcher = watcherClient.NewPurchaseWatcher()
	PurchaseWatcher.StartRoutine()

	DCAWatcher = watcherClient.NewDCAWatcher()
	DCAWatcher.StartRoutine()
}
