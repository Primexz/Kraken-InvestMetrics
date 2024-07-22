package kraken

import (
	"fmt"
	"strconv"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	jsonRequest "github.com/Primexz/Kraken-InvestMetrics/modules/http"
	"github.com/Primexz/go_kraken/rest"
	"github.com/shopspring/decimal"
)

type KrakenApi struct {
	api *rest.Kraken
}

type KrakenTickerInfo struct {
	Result map[string]struct {
		LastTrade []string `json:"c"`
	} `json:"result"`
}

func NewKraken() *KrakenApi {
	return &KrakenApi{
		api: rest.New(config.C.KrakenPublicKey, config.C.KrakenPrivateKey),
	}
}

func (k *KrakenApi) GetBtcOnKraken() (decimal.Decimal, error) {
	balance, err := k.api.GetAccountBalances()
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	return balance["XXBT"], nil
}

func (k *KrakenApi) GetPendingEuroOnKraken() (decimal.Decimal, error) {
	balance, err := k.api.GetAccountBalances()
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	return balance["ZEUR"], nil
}

func (k *KrakenApi) GetCachePayedToKraken() (float64, error) {
	ledgerInfo, err := k.api.GetLedgersInfo("deposit", 0, 0, false, "ZEUR")
	if err != nil {
		return 0, err
	}

	totalCachePayedToKraken := 0.0
	for _, value := range ledgerInfo {
		if value.Asset == "ZEUR" && value.LedgerType == "deposit" {
			totalCachePayedToKraken += value.Amount
		}
	}

	return totalCachePayedToKraken, nil
}

func (k *KrakenApi) GetAllBtcOrders(latestOnly bool) ([]rest.Ledger, error) {
	ledgers, err := k.api.GetLedgersInfo("", 0, 0, latestOnly, "XXBT")
	if err != nil {
		return nil, err
	}

	var orders []rest.Ledger

	for _, order := range ledgers {
		if order.Asset == "XXBT" && (order.LedgerType == "spend" || order.LedgerType == "margin" || order.LedgerType == "withdrawal" || order.LedgerType == "rollover" || order.Amount < 0) {
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (k *KrakenApi) GetCurrentBtcPrice(unit string) (float64, error) {
	pair := fmt.Sprintf("XXBTZ%s", unit)

	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%s", pair)

	resp, err := jsonRequest.GetJSON[KrakenTickerInfo](url)
	if err != nil {
		return 0, fmt.Errorf("failed to get btc price: %v %v", resp, err)
	}

	tickerInfo, ok := resp.Result[pair]
	if !ok {
		return 0, fmt.Errorf("failed to find btc price: %v", resp)
	}

	priceFloat, err := strconv.ParseFloat(tickerInfo.LastTrade[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse btc price: %v", resp)
	}

	return priceFloat, nil
}
