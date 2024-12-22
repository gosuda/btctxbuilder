package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/script"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type TxInput struct {
	tx *wire.MsgTx
	*types.Vin

	Amount        btcutil.Amount
	Address       btcutil.Address
	RedeemScript  []byte
	WitnessScript []byte
}

type TxInputs []*TxInput

func (t *TxInputs) AddInput(c *client.Client, txid string, vout uint32, amount int64, address string) error {
	tx, err := c.GetTx(txid)
	if err != nil {
		return err
	}
	if vout >= uint32(len(tx.Vout)) {
		return fmt.Errorf("vout %d out of range", vout)
	}
	prev := &tx.Vout[vout]

	rawTx, err := c.GetRawTx(txid)
	if err != nil {
		return err
	}
	msgTx, err := client.DecodeRawTransaction(rawTx)
	if err != nil {
		return err
	}

	btcAmount := btcutil.Amount(amount)
	btcAddress, err := types.DecodeAddress(address, c.GetParams())
	if err != nil {
		return err
	}

	vin := &types.Vin{
		Txid:    txid,
		Vout:    vout,
		Prevout: prev,
	}

	*t = append(*t, &TxInput{
		tx:      msgTx,
		Vin:     vin,
		Amount:  btcAmount,
		Address: btcAddress,
	})
	return nil
}

func (t *TxInputs) AmountTotal() btcutil.Amount {
	var total btcutil.Amount
	for _, input := range *t {
		total += input.Amount
	}
	return total
}

func (t *TxInputs) ToWire() ([]*wire.TxIn, error) {
	var txIns []*wire.TxIn = make([]*wire.TxIn, 0, len(*t))
	for _, in := range *t {
		txHash, err := chainhash.NewHashFromStr(in.Txid)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(txHash, in.Vout)

		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIns = append(txIns, txIn)
	}

	return txIns, nil
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
