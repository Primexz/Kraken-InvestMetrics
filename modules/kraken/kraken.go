package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/Primexz/go_kraken/rest"
	"github.com/shopspring/decimal"
)

type KrakenApi struct {
	api *rest.Kraken
}

type KrakenSpread struct {
	Error  []interface{}          `json:"error"`
	Result map[string]interface{} `json:"result"`
	Last   int                    `json:"last"`
}

func NewKraken() *KrakenApi {
	return &KrakenApi{
		api: rest.New(config.KrakenPublicKey, config.KrakenPrivateKey),
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
	ledgerInfo, err := k.api.GetLedgersInfo("", 0, 0, "ZEUR")
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

func (k *KrakenApi) GetAllBtcOrders() ([]rest.Ledger, error) {
	ledgers, err := k.api.GetLedgersInfo("", 0, 0, "XXBT")
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

func (k *KrakenApi) GetCurrentBtcPriceEur(unit string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.kraken.com/0/public/Spread?pair=%s", unit))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result KrakenSpread
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	allEurPrices := result.Result[unit].([]interface{})
	latestEurPrices := allEurPrices[len(allEurPrices)-1].([]interface{})
	currentEurPrice := latestEurPrices[len(latestEurPrices)-1]

	parsedEurPrice, err := strconv.ParseFloat(currentEurPrice.(string), 32)
	if err != nil {
		return 0, err
	}

	return parsedEurPrice, nil
}
