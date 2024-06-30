package jsonRequest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetJSON[T any](url string) (*T, error) {
	// #nosec G107
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request status code not OK: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var target T
	err = json.Unmarshal(body, &target)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &target, nil
}
