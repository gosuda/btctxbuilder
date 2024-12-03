package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
)

type TxBuilder struct {
	version int
	client  *client.Client
	params  *chaincfg.Params

	inputs  TxInputs
	outputs TxOutputs

	msgTx  *wire.MsgTx
	packet *psbt.Packet
}

func NewTxBuilder(cfg *chaincfg.Params, client *client.Client) *TxBuilder {
	return &TxBuilder{
		version: wire.TxVersion,
		params:  cfg,
		client:  client,
	}
}

func (t *TxBuilder) Build() (*psbt.Packet, error) {
	if len(t.inputs) == 0 && len(t.outputs) == 0 {
		return nil, fmt.Errorf("PSBT packet must contain at least one input or output")
	}

	outpoints, nSequences, err := t.inputs.ToWire()
	if err != nil {
		return nil, err
	}
	outputs, err := t.outputs.ToWire()
	if err != nil {
		return nil, err
	}

	t.msgTx = wire.NewMsgTx(int32(t.version))
	for i, in := range outpoints {
		t.msgTx.AddTxIn(&wire.TxIn{
			PreviousOutPoint: *in,
			Sequence:         nSequences[i],
		})
	}
	for _, out := range outputs {
		t.msgTx.AddTxOut(out)
	}

	p, err := psbt.NewFromUnsignedTx(t.msgTx)
	if err != nil {
		return nil, err
	}

	return p, nil
}
