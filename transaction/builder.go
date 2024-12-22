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

func NewTxBuilder(client *client.Client) *TxBuilder {
	return &TxBuilder{
		version: wire.TxVersion,
		params:  client.Params,
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
	for i, _ := range packet.Inputs {
		vin := t.inputs[i]

		addrType, err := types.GetAddressType(vin.Address, t.params)
		if err != nil {
			return err
		}

		// Set the WitnessUtxo or NonWitnessUtxo based on the address type
		if addrType == types.P2WPKH || addrType == types.P2WSH || addrType == types.TAPROOT {
			// For SegWit and Taproot, use WitnessUtxo
			packet.Inputs[i].WitnessUtxo = &wire.TxOut{
				Value:    int64(vin.Amount),
				PkScript: utils.MustDecode(vin.Prevout.Scriptpubkey),
			}
		} else {
			// For non-SegWit, we need to set full transaction for NonWitnessUtxo
			rawTxBytes, err := t.client.GetRawTx(vin.Txid)
			if err != nil {
				return err
			}
			msgTx, err := client.DecodeRawTransaction([]byte(rawTxBytes))
			if err != nil {
				return err
			}
			packet.Inputs[i].NonWitnessUtxo = msgTx
		}
	}
	return nil
}
