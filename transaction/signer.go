package transaction

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func SignTx(psbtRaw, privateKey []byte) (*wire.MsgTx, error) {
	priv, pub := btcec.PrivKeyFromBytes(privateKey)

	packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtRaw), false)
	if err != nil {
		return nil, err
	}
	sigHashes := txscript.NewTxSigHashes(packet.UnsignedTx, PsbtPrevOutputFetcher(packet))

	for i, input := range packet.Inputs {
		var sig []byte
		var err error
		switch {
		// P2PK script
		case txscript.IsPayToPubKey(input.WitnessUtxo.PkScript):
			sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
				input.WitnessUtxo.PkScript, txscript.SigHashAll, priv)
			if err != nil {
				return nil, fmt.Errorf("failed to sign P2PK: %v", err)
			}
		// P2PKH script
		case txscript.IsPayToPubKeyHash(input.WitnessUtxo.PkScript):
			sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
				input.NonWitnessUtxo.TxOut[i].PkScript, txscript.SigHashAll, priv)
			if err != nil {
				return nil, err
			}
		// P2SH script
		case txscript.IsPayToScriptHash(input.WitnessUtxo.PkScript):
			sig, err = txscript.RawTxInSignature(packet.UnsignedTx, i,
				input.RedeemScript, txscript.SigHashAll, priv)
			if err != nil {
				return nil, err
			}
		// P2WPKH script
		case txscript.IsPayToWitnessPubKeyHash(input.WitnessUtxo.PkScript):
			sig, err = txscript.RawTxInWitnessSignature(packet.UnsignedTx, sigHashes, i,
				input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
		// P2WSH script
		case txscript.IsPayToWitnessScriptHash(input.WitnessUtxo.PkScript):
			sig, err = txscript.RawTxInWitnessSignature(packet.UnsignedTx, sigHashes, i,
				input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
		// Multisig script
		case IsMultiSigScript(input.WitnessUtxo.PkScript):

		// P2TR script
		case txscript.IsPayToTaproot(input.WitnessUtxo.PkScript):
			tapRootSig, err := txscript.TaprootWitnessSignature(packet.UnsignedTx, sigHashes, i, input.WitnessUtxo.Value, input.WitnessScript, txscript.SigHashAll, priv)
			if err != nil {
				return nil, err
			}
			_ = tapRootSig

		// OP_RETURN script
		case txscript.IsNullData(input.WitnessUtxo.PkScript):

		default:
			return nil, fmt.Errorf("unsupported script type")
		}

		u, err := psbt.NewUpdater(packet)
		if err != nil {
			return nil, err
		}

		success, err := u.Sign(i, sig, pub.SerializeCompressed(), nil, nil)
		if err != nil {
			return nil, err
		} else if success != psbt.SignSuccesful {
			return nil, fmt.Errorf("signing failed, code: %d", success)
		}
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
