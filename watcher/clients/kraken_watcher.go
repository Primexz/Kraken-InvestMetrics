package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/kraken"
	"github.com/sirupsen/logrus"
)

type KrakenWatcher struct {
	CacheToKraken float64
	BtcOnKraken   float64
	PendingFiat   float64

	log *logrus.Entry
	api *kraken.KrakenApi

	interval time.Duration
}

func NewKrakenWatcher() *KrakenWatcher {
	return &KrakenWatcher{
		log: logrus.WithFields(logrus.Fields{
			"prefix": "kraken_watcher",
		}),
		api:      kraken.NewKraken(),
		interval: 30 * time.Second,
	}
}

func (kw *KrakenWatcher) UpdateData() {
	cacheToKraken, err := kw.api.GetCachePayedToKraken()
	if err != nil {
		kw.log.Error("Error getting Kraken data: ", err)
		return
	}

	btcOnKraken, err := kw.api.GetBtcOnKraken()
	if err != nil {
		kw.log.Error("Error getting Kraken data: ", err)
		return
	}

	pendingFiat, err := kw.api.GetPendingEuroOnKraken()
	if err != nil {
		kw.log.Error("Error getting Kraken data: ", err)
		return
	}

	kw.CacheToKraken = cacheToKraken
	kw.BtcOnKraken, _ = btcOnKraken.Float64()
	kw.PendingFiat, _ = pendingFiat.Float64()
}

func (kw *KrakenWatcher) GetInterval() time.Duration {
	return kw.interval
}