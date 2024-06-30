package util

import "github.com/sirupsen/logrus"

func LoggerWithPrefix(prefix string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"prefix": prefix,
	})
}
