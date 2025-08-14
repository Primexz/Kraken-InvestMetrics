package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/modules/lnd"
	"github.com/Primexz/Kraken-InvestMetrics/util"
	"github.com/sirupsen/logrus"
)

type LNDWatcher struct {
	SatAmount int64

	log      *logrus.Entry
	interval time.Duration
}

func NewLNDWatcher() *LNDWatcher {
	return &LNDWatcher{
		log:      util.LoggerWithPrefix("lnd_watcher"),
		interval: 5 * time.Minute,
	}
}

func (lndw *LNDWatcher) UpdateData() {
	rpcAddr := config.C.LndRpcAddress
	if rpcAddr == "" {
		return
	}

	client, err := lnd.NewLndClient()
	if err != nil {
		lndw.log.WithError(err).Error("Failed to create LND client")
		return
	}

	sat, err := client.GetTotalBalance()
	if err != nil {
		lndw.log.WithError(err).Error("Failed to get total balance")
		return
	}

	lndw.SatAmount = sat
}

func (lndw *LNDWatcher) GetInterval() time.Duration {
	return lndw.interval
}
