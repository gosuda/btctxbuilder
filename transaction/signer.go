package transaction

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func SignTx(net types.Network, psbtRaw, privateKey []byte) (*wire.MsgTx, error) {
	chain := types.GetParams(net)

	priv, pub := btcec.PrivKeyFromBytes(privateKey)
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtRaw), false)
	if err != nil {
		return nil, err
	}

	// sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, PsbtPrevOutputFetcher(packet))

	updater, err := psbt.NewUpdater(packet)
	if err != nil {
		return nil, err
	}

	for i, input := range packet.Inputs {
		var sig []byte
		scriptClass, addresses, nRequired, err := txscript.ExtractPkScriptAddrs(input.WitnessUtxo.PkScript, chain)
		if err != nil {
			return nil, err
		}

		switch scriptClass {
		case txscript.PubKeyTy: // P2PK
			err = signInputP2PK(updater, i, input.WitnessUtxo.PkScript, priv, txscript.SigHashAll)
		case txscript.PubKeyHashTy: // P2PKH
			err = signInputP2PKH(updater, i, input.NonWitnessUtxo, input.WitnessUtxo.PkScript, priv, txscript.SigHashAll)
		case txscript.ScriptHashTy: // P2SH

		case txscript.WitnessV0PubKeyHashTy: // P2WPKH

		case txscript.WitnessV0ScriptHashTy: // P2WSH

		case txscript.MultiSigTy: // Multisig

		case txscript.NullDataTy: // OP_RETURN

		case txscript.WitnessV1TaprootTy: // P2TR

		default:
			return nil, fmt.Errorf("unsupported script type")
		}
		if err != nil {
			return nil, err
		}

		success, err := updater.Sign(i, sig, pub.SerializeCompressed(), nil, nil)
		if err != nil {
			return nil, err
		} else if success != psbt.SignSuccesful {
			return nil, fmt.Errorf("signing failed, code: %d", success)
		}

		// switch {
		// // P2PK script
		// case txscript.IsPayToPubKey(input.WitnessUtxo.PkScript):
		// 	sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
		// 		input.WitnessUtxo.PkScript, txscript.SigHashAll, priv)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("failed to sign P2PK: %v", err)
		// 	}
		// // P2PKH script
		// case txscript.IsPayToPubKeyHash(input.WitnessUtxo.PkScript):
		// 	sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
		// 		input.NonWitnessUtxo.TxOut[i].PkScript, txscript.SigHashAll, priv)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// // P2SH script
		// case txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript):
		// 	sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
		// 		input.RedeemScript, txscript.SigHashAll, priv)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// // P2WPKH script
		// case txscript.IsPayToWitnessPubKeyHash(input.WitnessUtxo.PkScript):
		// 	sig, err = txscript.RawTxInWitnessSignature(packet.UnsignedTx, sigHashes, i,
		// 		input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
		// // P2WSH script
		// case txscript.IsPayToWitnessScriptHash(input.WitnessUtxo.PkScript):
		// 	sig, err = txscript.RawTxInWitnessSignature(packet.UnsignedTx, sigHashes, i,
		// 		input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
		// // Multisig script
		// case IsMultiSigScript(input.WitnessUtxo.PkScript):

		// // P2TR script
		// case txscript.IsPayToTaproot(input.WitnessUtxo.PkScript):
		// 	tapRootSig, err := txscript.TaprootWitnessSignature(packet.UnsignedTx, sigHashes, i, input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	_ = tapRootSig

		// // OP_RETURN script
		// case txscript.IsNullData(input.WitnessUtxo.PkScript):

		// default:
		// 	return nil, fmt.Errorf("unsupported script type")
		// }

	}

	err = psbt.MaybeFinalizeAll(packet)
	if err != nil {
		return nil, err
	}

	signedRawTx, err := psbt.Extract(packet)
	if err != nil {
		return nil, err
	}

	return signedRawTx, nil
}

func signInputP2PK(updater *psbt.Updater, i int, prevPkScript []byte, privKey *btcec.PrivateKey, hashType txscript.SigHashType) error {
	signature, err := txscript.RawTxInSignature(updater.Upsbt.UnsignedTx, i, prevPkScript, hashType, privKey)
	if err != nil {
		return err
	}

	if signOutcome, err := updater.Sign(i, signature, privKey.PubKey().SerializeCompressed(), nil, nil); err != nil {
		return err
	} else if signOutcome != psbt.SignSuccesful {
		return fmt.Errorf("signing failed, code: %d", signOutcome)
	}
	return nil
}

func signInputP2PKH(updater *psbt.Updater, i int, nonWitnessUtxo *wire.MsgTx, prevPkScript []byte, privKey *btcec.PrivateKey, hashType txscript.SigHashType) error {
	if err := updater.AddInNonWitnessUtxo(nonWitnessUtxo, i); err != nil {
		return err
	}

	signature, err := txscript.RawTxInSignature(updater.Upsbt.UnsignedTx, i, prevPkScript, hashType, privKey)
	if err != nil {
		return err
	}

	if signOutcome, err := updater.Sign(i, signature, privKey.PubKey().SerializeCompressed(), nil, nil); err != nil {
		return err
	} else if signOutcome != psbt.SignSuccesful {
		return fmt.Errorf("signing failed, code: %d", signOutcome)
	}

	return nil
}

func signInputP2WPKH(updater *psbt.Updater, i int, prevPkScript []byte, amt int64, privKey *btcec.PrivateKey, hashType txscript.SigHashType) error {

	return nil
}

func handleP2SH(updater *psbt.Updater, redeemScript []byte, i int, prevPkScript []byte, privKey *btcec.PrivateKey, hashType txscript.SigHashType, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
		return err
	}
	return nil
	// return signInputP2PKH(updater, i, in, redeemScript, privKey, hashType)
}

func signInputMultisig(updater *psbt.Updater, i int, addresses []btcutil.Address, nRequired int, prevPkScript []byte, sigScript, prevScript []byte) error {

	return nil
}

// P2PKH   : NonWitnessUtxo, PartialSigs, SighashType, FinalScriptSig
// P2SH    : NonWitnessUtxo, PartialSigs, SighashType, RedeemScript, FinalScriptSig
// P2WPKH  : WitnessUtxo, PartialSigs, SighashType, FinalScriptWitness
// P2WSH   : WitnessUtxo, PartialSigs, SighashType, WitnessScript, FinalScriptWitness
// Taproot ( Key Spend ) : WitnessUtxo, TaprootKeySpendSig, TaprootInternalKey
// Taproot ( Script Spend ) : WitnessUtxo, TaprootScriptSpendSig, TaprootLeafScript, TaprootMerkleRoot

func IsMultiSigScript(script []byte) bool {
	ok, _ := txscript.IsMultisigScript(script)
	return ok
}
