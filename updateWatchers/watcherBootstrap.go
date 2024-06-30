package watcher

import (
	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	watcherClient "github.com/Primexz/Kraken-InvestMetrics/updateWatchers/watchers"
	"github.com/sirupsen/logrus"
)

var (
	XPubWatcher     *watcherClient.XPubWatcher
	KrakenWatcher   *watcherClient.KrakenWatcher
	PurchaseWatcher *watcherClient.PurchaseWatcher
	DCAWatcher      *watcherClient.DCAWatcher
	UtxoWatcher     *watcherClient.UtxoWatcher

	log = logrus.WithFields(logrus.Fields{
		"prefix": "dca_bot_watcher",
	})
)

func BootstrapWatchers() {
	log.Info("Boostrapping watchers..")

	if xPub.IsXPub() {
		XPubWatcher = watcherClient.NewXPubWatcher()
		XPubWatcher.StartRoutine()

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
