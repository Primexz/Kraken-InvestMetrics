package config

import (
	"os"

	"github.com/primexz/KrakenDCA/logger"
)

var (
	KrakenPublicKey           string
	KrakenPrivateKey          string
	BitcoinAddress            string
	TimescaleConnectionString string
	BitcoinAccount            int
	BitcoinAddressGapLimit    int
	DCABotMetricUrl           string
	ElectrumServerAddress     string

	log *logger.Logger
)

func init() {
	log = logger.NewLogger("config")
	loadConfiguration()
}

func loadConfiguration() {
	log.Info("Loading configuration..")

	KrakenPublicKey = loadRequiredEnvVariable("KRAKEN_PUBLIC_KEY")
	KrakenPrivateKey = loadRequiredEnvVariable("KRAKEN_PRIVATE_KEY")
	BitcoinAddress = loadRequiredEnvVariable("INVEST_EXPORTER_BTC_ADDR")
	TimescaleConnectionString = loadRequiredEnvVariable("TIMESCALE_CONNECTION_STRING")
	DCABotMetricUrl = loadFallbackEnvVariable("DCA_BOT_METRIC_URL", "")
	ElectrumServerAddress = loadRequiredEnvVariable("ELECTRUM_SERVER_ADDRESS")

	BitcoinAccount = 0
	BitcoinAddressGapLimit = 20
}

func loadRequiredEnvVariable(envVar string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		log.Fatal("Required environment variable", envVar, "missing.")
	}

	return envData
}

func loadFallbackEnvVariable(envVar string, fallback string) string {
	envData := os.Getenv(envVar)

	if envData == "" {
		envData = fallback
	}

	return envData
}
