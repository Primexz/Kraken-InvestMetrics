package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/kraken"
	"github.com/Primexz/Kraken-InvestMetrics/util"
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
		log:      util.LoggerWithPrefix("kraken_watcher"),
		api:      kraken.NewKraken(),
		interval: 30 * time.Second,
	}
}

func (kw *KrakenWatcher) UpdateData() {
	cacheToKraken, err := kw.api.GetCachePayedToKraken()
	if err != nil {
		kw.log.Error("Error getting cache payed to kraken: ", err)
		return
	}

	btcOnKraken, err := kw.api.GetBtcOnKraken()
	if err != nil {
		kw.log.Error("Error getting btc on kraken: ", err)
		return
	}

	pendingFiat, err := kw.api.GetPendingEuroOnKraken()
	if err != nil {
		kw.log.Error("Error getting pending euro: ", err)
		return
	}

	kw.CacheToKraken = cacheToKraken
	kw.BtcOnKraken, _ = btcOnKraken.Float64()
	kw.PendingFiat, _ = pendingFiat.Float64()
}

func (kw *KrakenWatcher) GetInterval() time.Duration {
	return kw.interval
}
