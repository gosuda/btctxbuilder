package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
)

type TxBuilder struct {
	version int
	client  *client.Client
	params  *chaincfg.Params

	msgTx   *wire.MsgTx
	inputs  TxInputs
	outputs TxOutputs
}

func NewTxBuilder(cfg *chaincfg.Params, client *client.Client) *TxBuilder {
	version := wire.TxVersion
	return &TxBuilder{
		version: version,
		params:  cfg,
		client:  client,
	}
}

func (tx *TxBuilder) Build() (*psbt.Packet, error) {
	return nil, nil
}
