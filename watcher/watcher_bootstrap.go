package watcher

import (
	watcher_client "github.com/Primexz/Kraken-InvestMetrics/watcher/clients"
)

var (
	DCAWatcher      = watcher_client.NewDCAWatcher()
	KrakenWatcher   = watcher_client.NewKrakenWatcher()
	PurchaseWatcher = watcher_client.NewPurchaseWatcher()
	UtxoWatcher     = watcher_client.NewUtxoWatcher()
	XpubWatcher     = watcher_client.NewXPubWatcher()

	manager = NewWatcherManager()
)

func Load() {
	watchers := []Watcher{
		DCAWatcher,
		KrakenWatcher,
		PurchaseWatcher,
		UtxoWatcher,
		XpubWatcher,
	}

	for _, watcher := range watchers {
		manager.RegisterWatcher(watcher)
	}

	manager.StartAllWatchers()
}
