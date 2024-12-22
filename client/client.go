package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rabbitprincess/btctxbuilder/types"
)

const (
	// https://github.com/blockstream/esplora/blob/master/API.md
	ClientURL = "https://blockstream.info"
)

func NewClient(net types.Network) *Client {
	client := &Client{
		http:   http.DefaultClient,
		params: types.GetParams(net),
	}

	switch net {
	case types.BTC:
		client.params = &chaincfg.MainNetParams
		client.url = ClientURL + "/api"
	case types.BTC_Testnet3:
		client.params = &chaincfg.RegressionNetParams
		client.url = ClientURL + "/testnet/api"
	case types.BTC_Signet:
		client.params = &chaincfg.SigNetParams
		client.url = ClientURL + "/signet/api"
	}

	return client
}

type Client struct {
	params *chaincfg.Params
	url    string

	http *http.Client
}

func (t *Client) GetParams() *chaincfg.Params {
	return t.params
}

func (t *Client) Close() {
	if t.http == nil {
		return
	}
	t.http = nil
}

func RequestGet[T any](client *Client, endpoint string) (T, error) {
	fullURL := fmt.Sprintf("%s%s", client.url, endpoint)
	resp, err := client.http.Get(fullURL)
	if err != nil {
		return *new(T), fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(T), fmt.Errorf("failed to read response body: %w", err)
	}

	var result T
	if _, ok := any(result).(string); ok {
		return any(string(body)).(T), nil
	}
	if json.Valid(body) {
		if err := json.Unmarshal(body, &result); err != nil {
			return *new(T), fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return result, nil
	}

	// Non-JSON response, only supported for strings
	return *new(T), fmt.Errorf("non-JSON response cannot be parsed into %T", result)
}

func RequestPost[T any](client *Client, endpoint string, payload interface{}) (T, error) {
	fullURL := fmt.Sprintf("%s%s", client.url, endpoint)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return *new(T), fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.http.Post(fullURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return *new(T), fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(T), fmt.Errorf("failed to read response body: %w", err)
	}

	var result T
	if _, ok := any(result).(string); ok {
		return any(string(body)).(T), nil
	}
	if json.Valid(body) {
		if err := json.Unmarshal(body, &result); err != nil {
			return *new(T), fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return result, nil
	}

	// Non-JSON response, only supported for strings
	return *new(T), fmt.Errorf("non-JSON response cannot be parsed into %T", result)
}
