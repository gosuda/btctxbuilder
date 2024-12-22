package transaction

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/script"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type TxInput struct {
	NonWitnessUtxo *wire.MsgTx
	WitnessUtxo    *wire.TxOut
}

type TxInputs []*types.Vin

func (t *TxInputs) AddInputTransfer(txid string, vout uint32, address string, amount int64) error {
	vin := &types.Vin{
		Txid:    txid,
		Vout:    vout,
		Amount:  btcutil.Amount(amount),
		Address: address,
	}
	*t = append(*t, vin)
	return nil
}

func (t *TxInputs) AddInput(vin *types.Vin, address string, amount int64) error {

	vin.Address = address
	vin.Amount = btcutil.Amount(amount)
	*t = append(*t, vin)
	return nil
}

func (t *TxInputs) AmountTotal() btcutil.Amount {
	var total btcutil.Amount
	for _, vin := range *t {
		total += vin.Amount
	}
	return total
}

func (t *TxInputs) ToWire() ([]*wire.OutPoint, []uint32, error) {
	outpoints := make([]*wire.OutPoint, 0, len(*t))
	nSequences := make([]uint32, 0, len(*t))
	for _, in := range *t {
		txHash, err := chainhash.NewHashFromStr(in.Txid)
		if err != nil {
			return nil, nil, err
		}
		witness := make([][]byte, 0, len(in.Witness))
		for _, w := range in.Witness {
			witness = append(witness, []byte(w))
		}

		outpoints = append(outpoints, wire.NewOutPoint(txHash, in.Vout))
		nSequences = append(nSequences, wire.MaxTxInSequenceNum)
	}
	return outpoints, nSequences, nil
}

type TxOutput struct {
	Address btcutil.Address
	Amount  btcutil.Amount
}

type TxOutputs []*TxOutput

func (t *TxOutputs) AddOutputTransfer(params *chaincfg.Params, addr string, amount int64) error {
	rawAddr, err := types.DecodeAddress(addr, params)
	if err != nil {
		return err
	}
	vout := &TxOutput{
		Address: rawAddr,
		Amount:  btcutil.Amount(amount),
	}
	*t = append(*t, vout)
	return nil
}

func (t *TxOutputs) AmountTotal() btcutil.Amount {
	var total btcutil.Amount
	for _, vout := range *t {
		total += vout.Amount
	}
	return total
}

func (t *TxOutputs) ToWire() ([]*wire.TxOut, error) {
	txOuts := make([]*wire.TxOut, 0, len(*t))
	for _, out := range *t {
		pkScript, err := script.EncodeTransferScript(out.Address)
		if err != nil {
			return nil, err
		}
		txOuts = append(txOuts, wire.NewTxOut(int64(out.Amount), pkScript))
	}
	return txOuts, nil
}
