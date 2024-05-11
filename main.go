package main

import (
	"fmt"
	"runtime"

	"github.com/Primexz/Kraken-InvestMetrics/timescale"
	watcher "github.com/Primexz/Kraken-InvestMetrics/updateWatchers"

	"github.com/primexz/KrakenDCA/logger"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.NewLogger("main")
}

func main() {
	log.Info(fmt.Sprintf("Kraken Invest Metrics 🐙 %s, commit %s, built at %s (%s [%s, %s])", version, commit, date, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	watcher.BootstrapWatchers()
	go timescale.StartMetricRecorder()

	select {}
}
