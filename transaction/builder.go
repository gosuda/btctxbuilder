package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
)

type TxBuilder struct {
	version int32
	params  *chaincfg.Params
}

func NewTxBuilder(cfg *chaincfg.Params) *TxBuilder {
	version := int32(2)
	return &TxBuilder{
		version: version,
		params:  cfg,
	}
}

func (tx *TxBuilder) Build() (*psbt.Packet, error) {
	return nil, nil
}
