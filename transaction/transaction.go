package transaction

import (
	"fmt"
	"sort"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func NewTransferTx(net types.Network, fromAddress string, toAddress map[string]int64, changeAddress string) (*psbt.Packet, error) {
	client := client.NewClient(net)
	builder := NewTxBuilder(types.GetParams(net), client)

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
		builder.inputs.AddInputTransfer(utxo.Txid, utxo.Vout, fromAddress, utxo.Value)
	}

	// create outputs
	for address, amount := range toAddress {
		builder.outputs.AddOutput(address, amount)
	}

	// fund outputs
	if changeAddress == "" {
		changeAddress = fromAddress
	}
	err = builder.FundRawTransaction(changeAddress)
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func SelectUtxo(utxos []*types.Utxo, amount int64) (selected []*types.Utxo, totalAmount int64, err error) {
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value < utxos[j].Value
	})

	var total int64
	var selectedUtxos []*types.Utxo
	for _, utxo := range utxos {
		total += utxo.Value
		selectedUtxos = append(selectedUtxos, utxo)
		if total >= amount {
			break
		}
	}

	if total < amount {
		return nil, 0, fmt.Errorf("insufficient balance | total : %v | to amount : %v", total, amount)
	}

	return selectedUtxos, total, nil
}
