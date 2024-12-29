package transaction

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/script"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type TxInput struct {
	txid string
	vout uint32

	tx       *wire.MsgTx
	prevVout *wire.TxOut
	// *types.Vin

	Amount   btcutil.Amount
	Address  btcutil.Address
	AddrType types.AddrType

	// PkScript      []byte
	// RedeemScript  []byte
	// WitnessScript []byte
}

type TxInputs []*TxInput

func (t *TxInputs) AddInput(c *client.Client, txid string, vout uint32, amount int64, address string) error {
	// tx, err := c.GetTx(txid)
	// if err != nil {
	// 	return err
	// }
	// if vout >= uint32(len(tx.Vout)) {
	// 	return fmt.Errorf("vout %d out of range", vout)
	// }
	// prev := &tx.Vout[vout]

	rawTx, err := c.GetRawTx(txid)
	if err != nil {
		return err
	}
	msgTx, err := client.DecodeRawTransaction(rawTx)
	if err != nil {
		return err
	}
	prevVout := msgTx.TxOut[vout]

	btcAmount := btcutil.Amount(amount)
	btcAddress, _, err := types.DecodeAddress(address, c.GetParams())
	if err != nil {
		return err
	}
	AddrType := types.GetAddressType(btcAddress)

	// vin := &types.Vin{
	// 	Txid:    txid,
	// 	Vout:    vout,
	// 	Prevout: prev,
	// }

	*t = append(*t, &TxInput{
		txid:     txid,
		vout:     vout,
		tx:       msgTx,
		prevVout: prevVout,
		// Vin:      vin,
		Amount:   btcAmount,
		Address:  btcAddress,
		AddrType: AddrType,
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
		txHash, err := chainhash.NewHashFromStr(in.txid)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(txHash, in.vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIns = append(txIns, txIn)
	}

	return txIns, nil
}

type TxOutput struct {
	Amount   btcutil.Amount
	PkScript []byte
}

type TxOutputs []*TxOutput

func (t *TxOutputs) AddOutputTransfer(params *chaincfg.Params, addr string, amount int64) error {
	rawAddr, _, err := types.DecodeAddress(addr, params)
	if err != nil {
		return err
	}
	pkScript, err := script.EncodeTransferScript(rawAddr)
	if err != nil {
		return err
	}
	t.AddOutputPkScript(pkScript, amount)
	return nil
}

func (t *TxOutputs) AddOutputPkScript(pkScript []byte, amount int64) {
	vout := &TxOutput{
		Amount:   btcutil.Amount(amount),
		PkScript: pkScript,
	}
	*t = append(*t, vout)
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
		txOuts = append(txOuts, wire.NewTxOut(int64(out.Amount), out.PkScript))
	}
	return txOuts, nil
}
