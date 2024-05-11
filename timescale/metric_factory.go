package timescale

import (
	"context"
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/blockchain"
	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/kraken"
	watcher "github.com/Primexz/Kraken-InvestMetrics/updateWatchers"
	"github.com/Primexz/Kraken-InvestMetrics/xPub"
)

func StartMetricRecorder() {
	kraken := kraken.NewKraken()

	for {
		startTime := time.Now()

		btcOnKraken := watcher.KrakenWatcher.BtcOnKraken
		totalCache := watcher.KrakenWatcher.CacheToKraken
		pendingFiat := watcher.KrakenWatcher.PendingFiat

		btcEurPrice, err := kraken.GetCurrentBtcPriceEur("XXBTZEUR")
		if err != nil {
			LogMetricError(err)
			continue
		}

		btcUsdPrice, err := kraken.GetCurrentBtcPriceEur("XXBTZUSD")
		if err != nil {
			LogMetricError(err)
			continue
		}

		var walletBtc float64

		if xPub.IsXPub() {
			walletBtc = watcher.XPubWatcher.SatAmount
		} else {
			if addressInfo, err := blockchain.GetAddressInfo(config.BitcoinAddress); err == nil {
				walletBtc = float64(addressInfo.ChainStats.FundedTxoSum)
			} else {
				LogMetricError(err)
				continue
			}
		}

		//convert satoshi to btc
		walletBtc = walletBtc / 100000000

		ret, err := ConnectionPool.Exec(context.Background(), "INSERT INTO investment_exporter (time, total_btc_on_kraken, total_cache_to_kraken, eur_on_kraken, btc_price_eur, btc_price_usd, btc_in_wallet, eur_in_wallet, total_scrape_time) VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8)",
			btcOnKraken, totalCache-pendingFiat, btcOnKraken*btcEurPrice, btcEurPrice, btcUsdPrice, walletBtc, walletBtc*btcEurPrice, float64(time.Since(startTime).Milliseconds()))

		if err != nil {
			log.Error("failed to insert metrics into timescale: ", err, ret)
		}

		Timeout()
	}

}

func LogMetricError(err error) {
	log.Error("failed to scrape metrics: ", err)
	Timeout()
}

func Timeout() {
	time.Sleep(30 * time.Second)
}
