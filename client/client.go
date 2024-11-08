package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg"
)

func NewClient(net Network) *Client {
	client := &Client{
		net:  net,
		http: http.DefaultClient,
	}

	client.url = "https://blockstream.info"
	switch net {
	case Mainnet:
		client.params = &chaincfg.MainNetParams
		client.url = client.url + "/api"
	case Regtest:
		client.params = &chaincfg.RegressionNetParams
		client.url = client.url + "/testnet/api"
	case Signet:
		client.params = &chaincfg.SigNetParams
		client.url = client.url + "/signet/api"
	}

	return client
}

type Client struct {
	net    Network
	params *chaincfg.Params
	url    string

	http *http.Client
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
	if json.Valid(body) {
		if err := json.Unmarshal(body, &result); err != nil {
			return *new(T), fmt.Errorf("failed to unmarshal response: %w", err)
		}
	} else {
		// Non-JSON response, only supported for strings
		if _, ok := any(result).(string); ok {
			result = any(string(body)).(T)
		} else {
			return *new(T), fmt.Errorf("non-JSON response cannot be parsed into %T", result)
		}
	}
	return result, nil
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
	if json.Valid(body) {
		if err := json.Unmarshal(body, &result); err != nil {
			return *new(T), fmt.Errorf("failed to unmarshal response: %w", err)
		}
	} else {
		// Non-JSON response, only supported for strings
		if _, ok := any(result).(string); ok {
			result = any(string(body)).(T)
		} else {
			return *new(T), fmt.Errorf("non-JSON response cannot be parsed into %T", result)
		}
	}
	return result, nil
}
