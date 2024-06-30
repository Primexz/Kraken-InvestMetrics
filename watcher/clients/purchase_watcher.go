package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/kraken"
	"github.com/Primexz/Kraken-InvestMetrics/modules/timescale"
	"github.com/sirupsen/logrus"
)

type PurchaseWatcher struct {
	log *logrus.Entry
	api *kraken.KrakenApi

	interval time.Duration
}

func NewPurchaseWatcher() *PurchaseWatcher {
	return &PurchaseWatcher{
		log: logrus.WithFields(logrus.Fields{
			"prefix": "purchase_watcher",
		}),
		api:      kraken.NewKraken(),
		interval: 5 * time.Minute,
	}
}

func (pw *PurchaseWatcher) UpdateData() {
	//if the row count is greater than 50, then we only want to get the last 50 purchases
	purchases, err := pw.api.GetAllBtcOrders(pw.getPurchasesRowCount() > 50)
	if err != nil {
		pw.log.Error(err)
		return
	}

	for _, p := range purchases {
		time := time.Unix(int64(p.Time), 0)

		_, err := timescale.ConnectionPool.Exec(timescale.Context, "INSERT INTO purchases (refid, time, amount, fee) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING", p.RefID, time, p.Amount, p.Fee)
		if err != nil {
			pw.log.Error("failed to add purchase to db", err)
		}
	}
}

func (pw *PurchaseWatcher) GetInterval() time.Duration {
	return pw.interval
}

func (pw *PurchaseWatcher) getPurchasesRowCount() int {
	var count int
	err := timescale.ConnectionPool.QueryRow(timescale.Context, "SELECT COUNT(*) FROM purchases").Scan(&count)
	if err != nil {
		pw.log.Error("failed to get row count", err)
	}

	return count
}
