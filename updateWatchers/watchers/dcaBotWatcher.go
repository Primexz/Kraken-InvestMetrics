package watcherClient

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/primexz/KrakenDCA/logger"
)

type MetricResponse struct {
	NextOrder int64 `json:"nextOrder"`
}

type DCAWatcher struct {
	NextOrder time.Time

	log *logger.Logger
}

func NewDCAWatcher() *DCAWatcher {
	kw := &DCAWatcher{
		log: logger.NewLogger("dcawatcher"),
	}

	kw.UpdateData()

	return kw
}

func (dcaw *DCAWatcher) StartRoutine() {
	go func() {
		for {
			dcaw.UpdateData()
			time.Sleep(30 * time.Second)
		}
	}()
}

func (dcaw *DCAWatcher) UpdateData() {
	dcaw.log.Info("Updating DCA-Bot Watcher")

	url := config.C.DCABotMetricUrl
	if url == "" {
		return
	}

	// #nosec G107
	resp, err := http.Get(url)
	if err != nil {
		dcaw.log.Error(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dcaw.log.Error(err)
		return
	}

	metric := &MetricResponse{}
	err = json.Unmarshal(body, metric)
	if err != nil {
		dcaw.log.Error(err)
		return
	}

	dcaw.NextOrder = time.Unix(metric.NextOrder, 0)
	dcaw.log.Info("Next DCA order at ", dcaw.NextOrder)
}
