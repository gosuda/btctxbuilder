package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

type BuilderOpt func(*TxBuilder) error

func WithVersion(version int) BuilderOpt {
	return func(t *TxBuilder) error {
		t.version = version
		return nil
	}
}

func WithFundAddress(address string) BuilderOpt {
	return func(t *TxBuilder) error {
		t.fundAddress = address
		return nil
	}
}

type TxBuilder struct {
	version int
	client  *client.Client
	params  *chaincfg.Params

	inputs  TxInputs
	outputs TxOutputs

	fundAddress string

	msgTx  *wire.MsgTx
	packet *psbt.Packet
}

func NewTxBuilder(client *client.Client, opts ...BuilderOpt) *TxBuilder {
	builder := &TxBuilder{
		version: wire.TxVersion,
		params:  client.GetParams(),
		client:  client,
	}

	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func (t *TxBuilder) Build() (*psbt.Packet, error) {
	if len(t.inputs) == 0 && len(t.outputs) == 0 {
		return nil, fmt.Errorf("PSBT packet must contain at least one input or output")
	}

	txIns, err := t.inputs.ToWire()
	if err != nil {
		return nil, err
	}
	outputs, err := t.outputs.ToWire()
	if err != nil {
		return nil, err
	}

	t.msgTx = wire.NewMsgTx(int32(t.version))
	for _, in := range txIns {
		t.msgTx.AddTxIn(in)
	}
	for _, out := range outputs {
		t.msgTx.AddTxOut(out)
	}

	if t.fundAddress != "" {
		err = t.FundRawTransaction()
		if err != nil {
			return nil, err
		}
	}

	p, err := psbt.NewFromUnsignedTx(t.msgTx)
	if err != nil {
		return nil, err
	}

	err = t.decorateTxInputs(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (t *TxBuilder) decorateTxInputs(packet *psbt.Packet) error {
	for i := range packet.Inputs {
		vin := t.inputs[i]

		addrType := types.GetAddressType(vin.Address)
		if addrType == types.Invalid {
			return fmt.Errorf("invalid address type")
		}

		// Set the WitnessUtxo or NonWitnessUtxo based on the address type
		if addrType == types.P2WPKH || addrType == types.P2WSH || addrType == types.TAPROOT {
			// For SegWit and Taproot, use WitnessUtxo
			packet.Inputs[i].WitnessUtxo = &wire.TxOut{
				Value:    int64(vin.Amount),
				PkScript: utils.MustDecode(vin.Prevout.Scriptpubkey),
			}
		} else {
			packet.Inputs[i].NonWitnessUtxo = vin.tx
		}
		if vin.RedeemScript != nil {
			packet.Inputs[i].RedeemScript = vin.RedeemScript
		}
		if vin.WitnessScript != nil {
			packet.Inputs[i].WitnessScript = vin.WitnessScript
		}

	}
	return nil
}
