package watcherClient

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/xPub"
	"github.com/primexz/KrakenDCA/logger"
)

type UtxoWatcher struct {
	UtxoMap map[string]float64
	xPub    *xPub.XPub

	log *logger.Logger
}

func NewUtxoWatcher() *UtxoWatcher {
	u := &UtxoWatcher{
		xPub: xPub.NewXPub(),
		log:  logger.NewLogger("utxow"),
	}

	u.UpdateUtxoData()

	return u
}

func (u *UtxoWatcher) StartRoutine() {
	go func() {
		for {
			time.Sleep(30 * time.Minute)

			u.UpdateUtxoData()
		}
	}()
}

func (u *UtxoWatcher) UpdateUtxoData() {
	u.log.Info("Updating UTXO Watcher")

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
