package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/rabbitprincess/btctxbuilder/script"
	"github.com/rabbitprincess/btctxbuilder/types"
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
	feeAmount, err := EstimateTxFee(feeRate, t.inputs, t.msgTx.TxOut, changeAddressBTC)
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

func EstimateTxFee(feeRate float64, ins TxInputs, outs []*wire.TxOut, changeAddress btcutil.Address) (btcutil.Amount, error) {
	feeRatePerKb := btcutil.Amount(feeRate) * 1000
	vSize, err := EstimateTxVirtualSize(ins, outs, changeAddress)
	if err != nil {
		return 0, err
	}
	estimateFee := txrules.FeeForSerializeSize(feeRatePerKb, vSize)
	return estimateFee, nil
}

func EstimateTxVirtualSize(ins TxInputs, outs []*wire.TxOut, changeAddress btcutil.Address) (vSize int, err error) {
	// TODO : Add support for p2sh, p2wsh
	var nested, p2wpkh, p2tr, p2pkh int
	for _, in := range ins {
		switch types.GetAddressType(in.Address) {
		case types.P2PKH:
			p2pkh++
		case types.P2WPKH:
			p2wpkh++
		case types.P2WPKH_NESTED, types.P2WSH_NESTED:
			nested++
		case types.P2TR:
			p2tr++
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
