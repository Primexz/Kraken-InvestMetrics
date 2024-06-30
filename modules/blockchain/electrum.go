package blockchain

import (
	"context"
	"time"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/checksum0/go-electrum/electrum"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"prefix": "electrum",
	})

	Client *electrum.Client
	Ctx    = context.TODO()
)

func init() {
	electrum, err := electrum.NewClientTCP(Ctx, config.C.ElectrumServerAddress)
	if err != nil {
		log.Fatal("failed to connect to electrum server", err)
	}

	Client = electrum

	serverVersion, protocolVersion, err := Client.ServerVersion(Ctx)
	if err != nil {
		log.Fatal("failed to get server version", err)
	}

	log.Info("connected to electrum server", "server_version", serverVersion, "protocol_version", protocolVersion)

	banner, _ := Client.ServerBanner(Ctx)
	log.Info("server banner", banner)

	ping()
}

func GetBalance(address string) (float64, error) {

	scriptHash, err := electrum.AddressToElectrumScriptHash(address)
	if err != nil {
		log.Error("Failed to convert address to script hash:", err)
		return 0.0, err
	}

	balance, err := Client.GetBalance(Ctx, scriptHash)
	if err != nil {
		log.Error("Failed to resolve address:", err)
		return 0.0, err
	}

	return balance.Confirmed, nil
}

func ping() {
	go func() {
		for {
			if err := Client.Ping(Ctx); err != nil {
				log.Fatal(err)
			}

			time.Sleep(60 * time.Second)
		}
	}()
}
