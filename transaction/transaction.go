package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func NewTransferTx(c *client.Client, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) (*psbt.Packet, error) {
	var err error
	builder := NewTxBuilder(c)
	builder.FromAddress = fromAddress

	// fund fee outputs
	if fundAddress == "" {
		builder.FundAddress = fromAddress
	} else {
		builder.FundAddress = fundAddress
	}

	// estimate fee
	fees, err := c.FeeEstimate()
	if err != nil {
		return nil, err
	}
	builder.FeeRate = fees["1"]

	// create outputs
	for address, amount := range toAddress {
		if err = builder.Outputs.AddOutputTransfer(c.GetParams(), address, amount); err != nil {
			return nil, err
		}
	}
	toTotal := builder.Outputs.AmountTotal()

	// get utxo
	if len(utxos) == 0 { // if no utxos provided, get utxos from client
		utxos, err = builder.Client.GetUTXO(fromAddress)
		if err != nil {
			return nil, err
		}
	}

	// select utxo
	selected, unselected, err := SelectUtxo(utxos, int64(toTotal))
	if err != nil {
		return nil, err
	}
	// add inputs
	for _, utxo := range selected {
		if err = builder.Inputs.AddInput(c, utxo.Txid, utxo.Vout, utxo.Value, fromAddress); err != nil {
			return nil, err
		}
	}
	// unspent utxos
	builder.Utxos = unselected

	// build psbt from inputs and outputs
	return builder.Build()
}

func NewRunestoneEdictTx(c *client.Client, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) {

}
