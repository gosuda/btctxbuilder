package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type TxBuilder struct {
	txType  types.TxType
	version int

	client *client.Client
	params *chaincfg.Params
}

func NewTxBuilder(txType types.TxType, cfg *chaincfg.Params, client *client.Client) *TxBuilder {
	version := wire.TxVersion
	return &TxBuilder{
		txType:  txType,
		version: version,
		params:  cfg,
		client:  client,
	}
}

func (tx *TxBuilder) Build() (*psbt.Packet, error) {
	return nil, nil
}
