package client

import (
	"fmt"
)

func (c *Client) BestBlockHeight() (uint64, error) {
	return RequestGet[uint64](c, "/blocks/tip/height")
}

func (c *Client) BestBlockHash() (string, error) {
	return RequestGet[string](c, "/blocks/tip/hash")
}

func (c *Client) GetBlock(hash string) (*Block, error) {
	return RequestGet[*Block](c, fmt.Sprintf("/block/%s", hash))
}

func (c *Client) GetBlockTx(hash string, offset int) ([]*Transaction, error) {
	if offset > 0 {
		return RequestGet[[]*Transaction](c, fmt.Sprintf("/block/%s/txs/%d", hash, offset))
	}
	return RequestGet[[]*Transaction](c, fmt.Sprintf("/block/%s/txs", hash))
}

func (c *Client) GetTx(txid string) (*Transaction, error) {
	return RequestGet[*Transaction](c, fmt.Sprintf("/tx/%s", txid))
}

func (c *Client) GetRawTx(txid string) (string, error) {
	return RequestGet[string](c, fmt.Sprintf("/tx/%s/raw", txid))
}

func (c *Client) GetUTXO(address string) ([]*Utxo, error) {
	return RequestGet[[]*Utxo](c, fmt.Sprintf("/address/%s/utxo", address))
}

func (c *Client) BroadCastTx(rawTx string) (string, error) {
	return RequestPost[string](c, "/tx", rawTx)
}
