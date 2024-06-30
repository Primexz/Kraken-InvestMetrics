package watcherClient

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
}

func NewKrakenWatcher() *KrakenWatcher {
	kw := &KrakenWatcher{
		log: logrus.WithFields(logrus.Fields{
			"prefix": "kraken_watcher",
		}),
		api: kraken.NewKraken(),
	}

	kw.UpdateData()

	return kw
}

func (kw *KrakenWatcher) StartRoutine() {
	go func() {
		for {
			time.Sleep(2 * time.Minute)

			kw.UpdateData()
		}
	}()
}

func (kw *KrakenWatcher) UpdateData() {
	kw.log.Info("Updating Kraken Watcher")

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
