package lnd

import (
	"context"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type LndClient struct {
	client lnrpc.LightningClient
}

func NewLndClient() (*LndClient, error) {
	client, err := lndclient.NewBasicClient(config.C.LndRpcAddress, config.C.LndTlsCert, config.C.LndMacaroon, config.C.LndNetwork)
	if err != nil {
		log.WithError(err).Warn("Failed to create LND client")
		return nil, err
	}

	return &LndClient{
		client: client,
	}, nil

}

func (l *LndClient) GetTotalBalance() (int64, error) {
	walletBalance, err := l.client.WalletBalance(context.Background(), &lnrpc.WalletBalanceRequest{})
	if err != nil {
		log.WithError(err).Warn("Failed to get wallet balance")
		return 0, err
	}

	channelBalance, err := l.client.ChannelBalance(context.Background(), &lnrpc.ChannelBalanceRequest{})
	if err != nil {
		log.WithError(err).Warn("Failed to get channel balance")
		return 0, err
	}

	// #nosec G115
	return walletBalance.TotalBalance + int64(channelBalance.LocalBalance.Sat), nil
}
