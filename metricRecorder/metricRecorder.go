package metricRecorder

import (
	"context"
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/modules/kraken"
	"github.com/Primexz/Kraken-InvestMetrics/modules/timescale"
	"github.com/Primexz/Kraken-InvestMetrics/util"
	"github.com/Primexz/Kraken-InvestMetrics/watcher"
)

var (
	log = util.LoggerWithPrefix("metric_recorder")
)

func StartMetricRecorder() {
	log.Info("Initializing metric recorder")

	kraken := kraken.NewKraken()

	for {
		startTime := time.Now()

		btcOnKraken := watcher.KrakenWatcher.BtcOnKraken
		totalCache := watcher.KrakenWatcher.CacheToKraken
		pendingFiat := watcher.KrakenWatcher.PendingFiat

		btcEurPrice, err := kraken.GetCurrentBtcPriceEur("XXBTZEUR")
		if err != nil {
			logMetricError(err)
			continue
		}

		btcUsdPrice, err := kraken.GetCurrentBtcPriceEur("XXBTZUSD")
		if err != nil {
			logMetricError(err)
			continue
		}

		//convert satoshi to btc
		var walletBtc = watcher.WalletWatcher.SatAmount / 100_000_000
		var timeSpent = float64(time.Since(startTime).Milliseconds())

		ret, err := timescale.ConnectionPool.Exec(context.Background(), "INSERT INTO investment_exporter (time, total_btc_on_kraken, total_cache_to_kraken, btc_price_eur, btc_price_usd, btc_in_wallet, total_scrape_time, next_dca_order_time, pending_fiat) VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8)",
			btcOnKraken, totalCache-pendingFiat, btcEurPrice, btcUsdPrice, walletBtc, timeSpent, watcher.DCAWatcher.NextOrder, pendingFiat)

		if err != nil {
			log.Error("failed to insert metrics into timescale: ", err, ret)
		}

		for address, btc := range watcher.UtxoWatcher.UtxoMap {
			_, err := timescale.ConnectionPool.Exec(context.Background(), "INSERT INTO utxo_balances (address, btc) VALUES ($1, $2) ON CONFLICT (address) DO UPDATE SET btc = $2", address, btc)
			if err != nil {
				log.Error("failed to insert utxo metrics into timescale: ", err)
			}
		}

		timeout()
	}

}

func logMetricError(err error) {
	log.Error("failed to scrape metrics: ", err)
	timeout()
}

func timeout() {
	time.Sleep(30 * time.Second)
}
