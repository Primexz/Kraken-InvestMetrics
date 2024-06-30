package main

import (
	"runtime"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/Kraken-InvestMetrics/metricRecorder"
	"github.com/Primexz/Kraken-InvestMetrics/watcher"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat:  "2006/01/02 - 15:04:05",
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		SpacePadding:     45,
	})

	log.SetReportCaller(true)

	level, err := log.ParseLevel(config.C.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("Invalid log level")
	}

	log.SetLevel(level)
}

func main() {
	log.Infof("Kraken Invest Metrics 🐙 %s, commit %s, built at %s (%s [%s, %s])", version, commit, date, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	watcher.Load()
	metricRecorder.StartMetricRecorder()
}
