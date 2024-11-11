package transaction

import (
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcutil"
)

const (
	WitnessScaleFactor = 4
)

func (t *TxBuilder) FundRawTransaction(changeAddress string) error {
	totalInput := t.inputs.AmountTotal()
	totalOutput := t.outputs.AmountTotal()

	size := GetTxVirtualSize(btcutil.NewTx(t.msgTx))
	feeEstimate, err := t.client.FeeEstimate()
	if err != nil {
		return err
	}
	feeRate := feeEstimate["6"]
	fee := calculateFee(size, feeRate)

	change := totalInput - totalOutput - fee
	if change < 0 {
		return fmt.Errorf("insufficient funds, input: %d, output: %d, fee: %d", totalInput, totalOutput, fee)
	}

	// faucet remain vout
	// vout := t.outputs.AmountTotal()(faucetAddress)

	return nil
}

func GetTxVirtualSize(tx *btcutil.Tx) int64 {
	// vSize := (weight(tx) + 3) / 4
	//       := (((baseSize * 3) + totalSize) + 3) / 4
	// We add 3 here as a way to compute the ceiling of the prior arithmetic
	// to 4. The division by 4 creates a discount for wit witness data.
	return (GetTransactionWeight(tx) + (WitnessScaleFactor - 1)) / WitnessScaleFactor
}

func GetTransactionWeight(tx *btcutil.Tx) int64 {
	msgTx := tx.MsgTx()

	baseSize := msgTx.SerializeSizeStripped()
	totalSize := msgTx.SerializeSize()

	// (baseSize * 3) + totalSize
	return int64((baseSize * (WitnessScaleFactor - 1)) + totalSize)
}

func calculateFee(vsize int64, feeRate float64) int64 {
	return int64(math.Ceil(float64(vsize) * feeRate))
}
