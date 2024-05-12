package watcherClient

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/kraken"
	"github.com/Primexz/Kraken-InvestMetrics/modules/timescale"
	"github.com/primexz/KrakenDCA/logger"
)

type PurchaseWatcher struct {
	log *logger.Logger
	api *kraken.KrakenApi
}

func NewPurchaseWatcher() *PurchaseWatcher {
	pw := &PurchaseWatcher{
		log: logger.NewLogger("purchaseWatcher"),
		api: kraken.NewKraken(),
	}

	return pw
}

func (pw *PurchaseWatcher) StartRoutine() {
	go func() {
		for {
			pw.UpdateData()
			time.Sleep(5 * time.Minute)
		}
	}()
}

func (pw *PurchaseWatcher) UpdateData() {
	pw.log.Info("Updating Purchase Watcher")

	purchases, err := pw.api.GetAllBtcOrders()
	if err != nil {
		pw.log.Error(err)
		return
	}

	for _, p := range purchases {
		time := time.Unix(int64(p.Time), 0)
		timescale.ConnectionPool.Exec(timescale.Context, "INSERT INTO purchases (refid, time, amount, fee) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", p.RefID, time, p.Amount, p.Fee)
	}
}
