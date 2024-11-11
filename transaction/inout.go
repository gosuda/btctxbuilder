package transaction

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type TxInputs []*types.Vin

func (t TxInputs) AddInput(vin *types.Vin, address string, amount int64) error {
	if t == nil {
		t = make(TxInputs, 0)
	}

	vin.Address = address
	vin.Amount = amount
	t = append(t, vin)
	return nil
}

func (t TxInputs) AmountTotal() int64 {
	var total int64
	for _, vin := range t {
		total += vin.Amount
	}
	return total
}

type TxOutputs []*types.Vout

func (t TxOutputs) AddOutput(vout *types.Vout, address btcutil.Address, amount int64) error {
	if t == nil {
		t = make(TxOutputs, 0)
	}

	if vout == nil {

	}

	vout.Address = address.EncodeAddress()
	vout.Amount = amount
	t = append(t, vout)
	return nil
}

func (t *TxOutputs) AmountTotal() int64 {
	var total int64
	for _, vout := range *t {
		total += vout.Amount
	}
	return total
}
