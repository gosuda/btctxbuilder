package transaction

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
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

func (t TxInputs) ToWire() ([]*wire.OutPoint, []uint32, error) {
	outpoints := make([]*wire.OutPoint, 0, len(t))
	nSequences := make([]uint32, 0, len(t))
	for _, in := range t {
		txHash, err := chainhash.NewHashFromStr(in.Txid)
		if err != nil {
			return nil, nil, err
		}
		witness := make([][]byte, 0, len(in.Witness))
		for _, w := range in.Witness {
			witness = append(witness, []byte(w))
		}

		outpoints = append(outpoints, wire.NewOutPoint(txHash, in.Vout))
		nSequences = append(nSequences, uint32(in.Sequence))
	}
	return outpoints, nSequences, nil
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

func (t TxOutputs) AmountTotal() int64 {
	var total int64
	for _, vout := range t {
		total += vout.Amount
	}
	return total
}

func (t TxOutputs) ToWire() ([]*wire.TxOut, error) {
	txOuts := make([]*wire.TxOut, 0, len(t))
	for _, out := range t {
		pkScript, err := utils.Decode(out.Scriptpubkey)
		if err != nil {
			return nil, err
		}
		txOuts = append(txOuts, wire.NewTxOut(out.Amount, pkScript))
	}
	return txOuts, nil
}
