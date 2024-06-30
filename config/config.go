package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/primexz/KrakenDCA/logger"
)

type config struct {
	KrakenPublicKey           string `env:"KRAKEN_PUBLIC_KEY,required"`
	KrakenPrivateKey          string `env:"KRAKEN_PRIVATE_KEY,required"`
	BitcoinAddress            string `env:"INVEST_EXPORTER_BTC_ADDR,required"`
	TimescaleConnectionString string `env:"TIMESCALE_CONNECTION_STRING,required"`
	BitcoinAccount            int    `env:"INVEST_EXPORTER_BTC_ACCOUNT" envDefault:"0"`
	BitcoinAddressGapLimit    int    `env:"INVEST_EXPORTER_BTC_GAP_LIMIT" envDefault:"20"`
	DCABotMetricUrl           string `env:"DCA_BOT_METRIC_URL"`
	ElectrumServerAddress     string `env:"ELECTRUM_SERVER_ADDRESS,required"`
}

var (
	log *logger.Logger
	C   config
)

func init() {
	log = logger.NewLogger("config")
	loadConfiguration()
}

func loadConfiguration() {
	if config, err := env.ParseAs[config](); err == nil {
		C = config
	} else {
		log.Fatal(err)
	}
}
