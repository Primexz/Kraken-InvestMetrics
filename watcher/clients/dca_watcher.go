package watcher_client

import (
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	jsonRequest "github.com/Primexz/Kraken-InvestMetrics/modules/http"
	"github.com/sirupsen/logrus"
)

type MetricResponse struct {
	NextOrder int64 `json:"nextOrder"`
}

type DCAWatcher struct {
	interval time.Duration

	NextOrder time.Time
	log       *logrus.Entry
}

func NewDCAWatcher() *DCAWatcher {
	return &DCAWatcher{
		log: logrus.WithFields(logrus.Fields{
			"prefix": "dca_bot_watcher",
		}),
		interval: 30 * time.Second,
	}
}

func (dcaw *DCAWatcher) UpdateData() {
	url := config.C.DCABotMetricUrl
	if url == "" {
		return
	}

	resp, err := jsonRequest.GetJSON[MetricResponse](url)
	if err != nil {
		dcaw.log.Error(err)
		return
	}
	dcaw.NextOrder = time.Unix(resp.NextOrder, 0)
	dcaw.log.Info("Next DCA order at ", dcaw.NextOrder)
}

func (dcaw *DCAWatcher) GetInterval() time.Duration {
	return dcaw.interval
}
