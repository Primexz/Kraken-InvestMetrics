package timescale

import "github.com/sirupsen/logrus"

var (
	log = logrus.WithFields(logrus.Fields{
		"prefix": "timescale",
	})
)
