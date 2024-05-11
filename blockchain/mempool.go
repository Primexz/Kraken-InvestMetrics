package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DifficultyAdjustmentInfo struct {
	ProgressPercent       float64 `json:"progressPercent"`
	DifficultyChange      float64 `json:"difficultyChange"`
	EstimatedRetargetDate int64   `json:"estimatedRetargetDate"`
	RemainingBlocks       int     `json:"remainingBlocks"`
	RemainingTime         int     `json:"remainingTime"`
	PreviousRetarget      float64 `json:"previousRetarget"`
	PreviousTime          int     `json:"previousTime"`
	NextRetargetHeight    int     `json:"nextRetargetHeight"`
	TimeAvg               int     `json:"timeAvg"`
	AdjustedTimeAvg       int     `json:"adjustedTimeAvg"`
	TimeOffset            int     `json:"timeOffset"`
	ExpectedBlocks        float64 `json:"expectedBlocks"`
}

type AddressInfo struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoSum int `json:"funded_txo_sum"`
		SpentTxoSum  int `json:"spent_txo_sum"`
	} `json:"chain_stats"`
}

func GetAddressInfo(addr string) (AddressInfo, error) {
	resp, err := http.Get("https://mempool.space/api/address/" + addr)
	if err != nil {
		return AddressInfo{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressInfo{}, err
	}

	var result AddressInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return AddressInfo{}, fmt.Errorf("failed to unmarshal response: %w %s", err, string(body))
	}

	return result, nil
}
