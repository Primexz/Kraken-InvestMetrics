package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/modules/blockchain"
	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	"github.com/Primexz/Kraken-InvestMetrics/util"
	"github.com/sirupsen/logrus"
)

type WalletWatcher struct {
	SatAmount float64

	xPub *xPub.XPub

	log      *logrus.Entry
	interval time.Duration
}

func NewWalletWatcher() *WalletWatcher {
	return &WalletWatcher{
		xPub:     xPub.NewXPub(),
		log:      util.LoggerWithPrefix("wallet_watcher"),
		interval: 5 * time.Minute,
	}
}

func (ww *WalletWatcher) UpdateData() {

	if util.IsXPub() {
		if amount, err := ww.xPub.GetTotalSats(); err == nil {
			ww.SatAmount = amount
		} else {
			ww.log.Error("failed to get total bitcoin amount", err)
		}

	} else {
		if balance, err := blockchain.GetBalance(config.C.BitcoinAddress); err == nil {
			ww.SatAmount = balance
		} else {
			ww.log.Error("failed to get total bitcoin amount", err)
		}
	}
}

func (ww *WalletWatcher) GetInterval() time.Duration {
	return ww.interval
}
