package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/rabbitprincess/btctxbuilder/script"
)

const (
	WitnessScaleFactor = 4
)

func (t *TxBuilder) SufficentFunds() bool {
	inSum, _ := btcutil.NewAmount(0)
	for _, input := range t.inputs {
		inSum += input.Amount
	}

	outSum, _ := btcutil.NewAmount(0)
	for _, output := range t.outputs {
		outSum += output.Amount
	}

	change := inSum - outSum
	return change >= 0
}

func (t *TxBuilder) FundRawTransaction() error {
	changeAddressBTC, err := btcutil.DecodeAddress(t.fundAddress, t.params)
	if err != nil {
		return err
	}

	// calculate fee
	feeEstimate, err := t.client.FeeEstimate()
	if err != nil {
		return err
	}
	feeRate := feeEstimate["6"]
	feeAmount, err := EstimateTxFee(feeRate, t.msgTx.TxIn, t.msgTx.TxOut, changeAddressBTC)
	if err != nil {
		return err
	}

	// calculate change amount
	totalInput := t.inputs.AmountTotal()
	totalOutput := t.outputs.AmountTotal()
	change := totalInput - totalOutput - feeAmount
	if change < 0 {
		return fmt.Errorf("insufficient funds, input: %d, output: %d, fee: %d", totalInput, totalOutput, feeAmount)
	}

	// add change output
	if change > 0 {
		pkScript, err := script.EncodeTransferScript(changeAddressBTC)
		if err != nil {
			return err
		}
		changeTxOut := wire.NewTxOut(int64(change), pkScript)
		t.msgTx.TxOut = append(t.msgTx.TxOut, changeTxOut)
	}

	return nil
}

func EstimateTxFee(feeRate float64, ins []*wire.TxIn, outs []*wire.TxOut, changeAddress btcutil.Address) (btcutil.Amount, error) {
	feeRatePerKb := btcutil.Amount(feeRate)
	vSize, err := EstimateTxVirtualSize(ins, outs, changeAddress)
	if err != nil {
		return 0, err
	}
	estimateFee := txrules.FeeForSerializeSize(feeRatePerKb, vSize)
	return estimateFee, nil
}

func EstimateTxVirtualSize(ins []*wire.TxIn, outs []*wire.TxOut, changeAddress btcutil.Address) (vSize int, err error) {
	var nested, p2wpkh, p2tr, p2pkh int
	for _, in := range ins {
		// Check the script type for each input
		if len(in.Witness) > 0 {
			if len(in.Witness[0]) == 64 { // Assuming Schnorr signature size
				p2tr++
			} else {
				p2wpkh++
			}
		} else if len(in.SignatureScript) > 0 {
			if len(in.SignatureScript) == txsizes.RedeemNestedP2WPKHInputSize {
				nested++
			} else {
				p2pkh++
			}
		}
	}
	changeScriptSize, err := GetChangeScriptSize(changeAddress)
	if err != nil {
		return 0, err
	}

	vSize = txsizes.EstimateVirtualSize(p2pkh, p2tr, p2wpkh, nested, outs, changeScriptSize)
	return vSize, nil
}

func GetChangeScriptSize(changeAddress btcutil.Address) (int, error) {
	// Determine the script type and size
	switch changeAddress.(type) {
	case *btcutil.AddressPubKeyHash: // P2PKH
		return 25, nil // OP_DUP OP_HASH160 [20-byte HASH] OP_EQUALVERIFY OP_CHECKSIG
	case *btcutil.AddressScriptHash: // P2SH
		return 23, nil // OP_HASH160 [20-byte HASH] OP_EQUAL
	case *btcutil.AddressWitnessPubKeyHash: // P2WPKH
		return 22, nil // OP_0 [20-byte HASH]
	case *btcutil.AddressWitnessScriptHash: // P2WSH
		return 34, nil // OP_0 [32-byte HASH]
	case *btcutil.AddressTaproot: // P2TR
		return 34, nil // OP_1 [32-byte Taproot PubKey]
	default:
		return 0, fmt.Errorf("unsupported address type: %T", changeAddress)
	}
}
