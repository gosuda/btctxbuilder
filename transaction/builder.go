package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
)

type TxType string

const (
	Transfer TxType = "transfer"
	FeeBump  TxType = "feebump"

	Script TxType = "script"
)

type TxBuilder struct {
	txType  TxType
	version int

	client *client.Client
	params *chaincfg.Params
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
