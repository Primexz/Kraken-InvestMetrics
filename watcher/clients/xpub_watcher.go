package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	"github.com/sirupsen/logrus"
)

type XPubWatcher struct {
	SatAmount float64
	xPub      *xPub.XPub

	log      *logrus.Entry
	interval time.Duration
}

func NewXPubWatcher() *XPubWatcher {
	return &XPubWatcher{
		xPub: xPub.NewXPub(),
		log: logrus.WithFields(logrus.Fields{
			"prefix": "xpub_watcher",
		}),
		interval: 5 * time.Minute,
	}
}

func (xw *XPubWatcher) UpdateData() {
	amount, err := xw.xPub.GetTotalSats()
	if err != nil {
		xw.log.Error("failed to get total bitcoin amount", err)
		return
	}

	xw.SatAmount = amount
}

func (xw *XPubWatcher) GetInterval() time.Duration {
	return xw.interval
}
