package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/gosuda/btctxbuilder/types"
)

func NewTransferTx(params *chaincfg.Params, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string, fee float64) (*psbt.Packet, error) {
	var err error

	builder := NewTxBuilder(params)
	builder.FromAddress = fromAddress

	// fund fee outputs
	if fundAddress == "" {
		builder.FundAddress = fromAddress
	} else {
		builder.FundAddress = fundAddress
	}

	// estimate fee
	builder.FeeRate = fee

	// create outputs
	for address, amount := range toAddress {
		if err = builder.Outputs.AddOutputTransfer(params, address, amount); err != nil {
			return nil, err
		}
	}
	toTotal := builder.Outputs.AmountTotal()

	// select utxo
	selected, unselected, err := SelectUtxo(utxos, int64(toTotal))
	if err != nil {
		return nil, err
	}
	// add inputs
	for _, utxo := range selected {
		if err = builder.Inputs.AddInput(params, utxo.RawTx, utxo.Vout, utxo.Value, fromAddress); err != nil {
			return nil, err
		}
	}
	// unspent utxos
	builder.Utxos = unselected

	// build psbt from inputs and outputs
	return builder.Build()
}

func NewRunestoneEdictTx(params *chaincfg.Params, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) {

}
