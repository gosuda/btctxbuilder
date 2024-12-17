package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func NewTransferTx(net types.Network, fromAddress string, toAddress map[string]int64, changeAddress string) (*psbt.Packet, error) {
	c := client.NewClient(net)
	builder := NewTxBuilder(types.GetParams(net), c)

	var toTotal int64
	for _, amount := range toAddress {
		toTotal += amount
	}

	// get utxo
	utxos, err := builder.client.GetUTXO(fromAddress)
	if err != nil {
		return nil, err
	}

	// select utxo
	selected, _, err := SelectUtxo(utxos, toTotal)
	if err != nil {
		return nil, err
	}

	// create inputs
	for _, utxo := range selected {
		rawTxBytes, err := builder.client.GetRawTx(utxo.Txid)
		if err != nil {
			return nil, err
		}
		msgTx, err := client.DecodeRawTransaction([]byte(rawTxBytes))
		if err != nil {
			return nil, err
		}
		_ = msgTx
		// vout := msgTx.TxOut[utxo.Vout]
		builder.inputs.AddInputTransfer(utxo.Txid, utxo.Vout, fromAddress, utxo.Value)
	}

	// create outputs
	for address, amount := range toAddress {
		builder.outputs.AddOutput(address, amount)
	}

	// fund outputs
	if changeAddress == "" {
		builder.changeAddress = fromAddress
	} else {
		builder.changeAddress = changeAddress
	}

	return builder.Build()
}
