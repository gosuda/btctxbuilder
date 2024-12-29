package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func NewTransferTx(c *client.Client, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) (*psbt.Packet, error) {
	var err error
	builder := NewTxBuilder(c)

	var toTotal int64
	for _, amount := range toAddress {
		toTotal += amount
	}

	// if no utxos provided, get utxos from client
	if len(utxos) == 0 {
		utxos, err = builder.client.GetUTXO(fromAddress)
		if err != nil {
			return nil, err
		}
	}

	// select utxo
	selected, _, err := SelectUtxo(utxos, toTotal)
	if err != nil {
		return nil, err
	}

	// create inputs
	for _, utxo := range selected {
		if err = builder.inputs.AddInput(c, utxo.Txid, utxo.Vout, utxo.Value, fromAddress); err != nil {
			return nil, err
		}
	}

	// create outputs
	for address, amount := range toAddress {
		if err = builder.outputs.AddOutputTransfer(c.GetParams(), address, amount); err != nil {
			return nil, err
		}
	}

	// fund fee outputs
	if fundAddress == "" {
		builder.fundAddress = fromAddress
	} else {
		builder.fundAddress = fundAddress
	}

	// build psbt from inputs and outputs
	return builder.Build()
}

func NewRunestoneEdictTx(c *client.Client, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) {

}
