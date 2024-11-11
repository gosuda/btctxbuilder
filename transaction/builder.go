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
	return &TxBuilder{
		version: wire.TxVersion,
		params:  cfg,
		client:  client,
	}
}

func (t *TxBuilder) Build() (*psbt.Packet, error) {
	outpoints, nSequences, err := t.inputs.ToWire()
	if err != nil {
		return nil, err
	}
	outputs, err := t.outputs.ToWire()
	if err != nil {
		return nil, err
	}
	return psbt.New(outpoints, outputs, int32(t.version), 0, nSequences)
}
