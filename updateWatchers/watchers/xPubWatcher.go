package watcherClient

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	"github.com/sirupsen/logrus"
)

type XPubWatcher struct {
	SatAmount float64
	xPub      *xPub.XPub

	log *logrus.Entry
}

func NewXPubWatcher() *XPubWatcher {
	xw := &XPubWatcher{
		xPub: xPub.NewXPub(),
		log: logrus.WithFields(logrus.Fields{
			"prefix": "xpub_watcher",
		}),
	}

	xw.UpdateCoinAmount()

	return xw
}

func (xw *XPubWatcher) StartRoutine() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)

			xw.UpdateCoinAmount()
		}
	}()
}

func (xw *XPubWatcher) UpdateCoinAmount() {
	xw.log.Info("Updating xPub Watcher")

	amount, err := xw.xPub.GetTotalSats()
	if err != nil {
		xw.log.Error("failed to get total bitcoin amount", err)
		return
	}

	xw.SatAmount = amount
}
