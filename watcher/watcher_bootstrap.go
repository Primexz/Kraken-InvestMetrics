package watcher

import (
	watcher_client "github.com/Primexz/Kraken-InvestMetrics/watcher/clients"
)

var (
	DCAWatcher      = watcher_client.NewDCAWatcher()
	KrakenWatcher   = watcher_client.NewKrakenWatcher()
	PurchaseWatcher = watcher_client.NewPurchaseWatcher()
	UtxoWatcher     = watcher_client.NewUtxoWatcher()
	WalletWatcher   = watcher_client.NewWalletWatcher()
	LNDWatcher      = watcher_client.NewLNDWatcher()

	manager = NewWatcherManager()
)

func Load() {
	watchers := []Watcher{
		DCAWatcher,
		KrakenWatcher,
		PurchaseWatcher,
		UtxoWatcher,
		WalletWatcher,
		LNDWatcher,
	}

	for _, watcher := range watchers {
		manager.RegisterWatcher(watcher)
	}

	manager.StartAllWatchers()
}
