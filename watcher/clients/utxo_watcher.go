package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	"github.com/sirupsen/logrus"
)

type UtxoWatcher struct {
	UtxoMap map[string]float64
	xPub    *xPub.XPub

	log      *logrus.Entry
	interval time.Duration
}

func NewUtxoWatcher() *UtxoWatcher {
	return &UtxoWatcher{
		xPub: xPub.NewXPub(),
		log: logrus.WithFields(logrus.Fields{
			"prefix": "utxo_watcher",
		}),
		interval: 5 * time.Minute,
	}
}

func (u *UtxoWatcher) UpdateData() {
	utxoMap, err := u.xPub.GetAddressSatMap()
	if err != nil {
		u.log.Error("failed to get utxo map", err)
		return
	}

	bitcoinMap := make(map[string]float64)

	for address, sat := range utxoMap {
		bitcoinMap[address] = sat / 100_000_000
	}

	u.UtxoMap = bitcoinMap
}

func (u *UtxoWatcher) GetInterval() time.Duration {
	return u.interval
}
