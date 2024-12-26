package client

import (
	"fmt"

	"github.com/rabbitprincess/btctxbuilder/types"
)

func (c *Client) BestBlockHeight() (uint64, error) {
	return RequestGet[uint64](c, "/blocks/tip/height")
}

func (c *Client) BestBlockHash() (string, error) {
	return RequestGet[string](c, "/blocks/tip/hash")
}

func (c *Client) GetBlock(hash string) (*types.Block, error) {
	return RequestGet[*types.Block](c, fmt.Sprintf("/block/%s", hash))
}

func (c *Client) GetBlockTx(hash string, offset int) ([]*types.Transaction, error) {
	if offset > 0 {
		return RequestGet[[]*types.Transaction](c, fmt.Sprintf("/block/%s/txs/%d", hash, offset))
	}
	return RequestGet[[]*types.Transaction](c, fmt.Sprintf("/block/%s/txs", hash))
}

func (c *Client) GetTx(txid string) (*types.Transaction, error) {
	return RequestGet[*types.Transaction](c, fmt.Sprintf("/tx/%s", txid))
}

func (c *Client) GetRawTx(txid string) (string, error) {
	return RequestGet[string](c, fmt.Sprintf("/tx/%s/raw", txid))
}

func (c *Client) GetAddress(address string) (*types.Address, error) {
	return RequestGet[*types.Address](c, fmt.Sprintf("/address/%s", address))
}

func (c *Client) GetUTXO(address string) ([]*types.Utxo, error) {
	return RequestGet[[]*types.Utxo](c, fmt.Sprintf("/address/%s/utxo", address))
}

func (c *Client) FeeEstimate() (types.FeeEstimate, error) {
	return RequestGet[types.FeeEstimate](c, "/fee-estimates")
}

func (c *Client) BroadcastTx(rawTx string) (string, error) {
	return RequestPost[string](c, "/tx", rawTx)
}
