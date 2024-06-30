package util

import (
	"strings"

	"github.com/Primexz/Kraken-InvestMetrics/config"
)

func IsXPub() bool {
	return strings.HasPrefix(config.C.BitcoinAddress, "xpub")
}

func Sum(arr []float64) float64 {
	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	return sum
}
