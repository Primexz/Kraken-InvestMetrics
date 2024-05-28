package watcher

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	watcherClient "github.com/Primexz/Kraken-InvestMetrics/updateWatchers/watchers"
	"github.com/primexz/KrakenDCA/logger"
)

var (
	XPubWatcher     *watcherClient.XPubWatcher
	KrakenWatcher   *watcherClient.KrakenWatcher
	PurchaseWatcher *watcherClient.PurchaseWatcher
	DCAWatcher      *watcherClient.DCAWatcher
	UtxoWatcher     *watcherClient.UtxoWatcher

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

		time.Sleep(5 * time.Second)

		UtxoWatcher = watcherClient.NewUtxoWatcher()
		UtxoWatcher.StartRoutine()
	}

	KrakenWatcher = watcherClient.NewKrakenWatcher()
	KrakenWatcher.StartRoutine()

	PurchaseWatcher = watcherClient.NewPurchaseWatcher()
	PurchaseWatcher.StartRoutine()

	DCAWatcher = watcherClient.NewDCAWatcher()
	DCAWatcher.StartRoutine()
}
