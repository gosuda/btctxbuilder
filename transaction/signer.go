package transaction

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func SignTx(net types.Network, psbtRaw, privateKey []byte) ([]byte, error) {
	chain := types.GetParams(net)

	priv, _ := btcec.PrivKeyFromBytes(privateKey)
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtRaw), false)
	if err != nil {
		return nil, err
	}

	err = psbt.InputsReadyToSign(packet)
	if err != nil {
		return nil, err
	}

	updater, err := psbt.NewUpdater(packet)
	if err != nil {
		return nil, err
	}

	prevOutputFetcher := PsbtPrevOutputFetcher(packet)

	for i, input := range packet.Inputs {

		// Extract previous transaction output information
		if input.WitnessUtxo == nil && input.NonWitnessUtxo == nil {
			return nil, fmt.Errorf("missing input UTXO information for input %d", i)
		}

		var prevOutValue int64
		var pkScript []byte
		if input.WitnessUtxo != nil {
			updater.AddInWitnessUtxo(input.WitnessUtxo, i)
			// prevOutValue = input.WitnessUtxo.Value
			pkScript = input.WitnessUtxo.PkScript
		} else if input.NonWitnessUtxo != nil {
			index := packet.UnsignedTx.TxIn[i].PreviousOutPoint.Index
			prevOut := input.NonWitnessUtxo.TxOut[index]
			// prevOutValue = prevOut.Value
			pkScript = prevOut.PkScript
		}
		_ = prevOutValue

		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(pkScript, chain)
		if err != nil {
			return nil, err
		}

		switch scriptClass {
		case txscript.WitnessV1TaprootTy: // P2TR
			err = signInputTaproot(updater, priv, i, pkScript, prevOutputFetcher)
		case txscript.PubKeyTy: // P2PK
			err = signInputP2PK(updater, i, pkScript, priv)
		case txscript.PubKeyHashTy: // P2PKH
			err = signInputP2PKH(updater, i, pkScript, priv)
		case txscript.ScriptHashTy: // P2SH
			panic("not supported yet")
		case txscript.WitnessV0PubKeyHashTy: // P2WPKH
			err = signInputP2WPKH(updater, i, pkScript, input.WitnessUtxo.Value, priv, prevOutputFetcher)
		case txscript.WitnessV0ScriptHashTy: // P2WSH
			panic("not supported yet")
		case txscript.MultiSigTy: // Multisig
			panic("not supported yet")
		case txscript.NullDataTy: // OP_RETURN

		default:
			return nil, fmt.Errorf("unsupported script type")
		}
		if err != nil {
			return nil, err
		}
		err = psbt.Finalize(packet, i)
		if err != nil {
			return nil, err
		}
	}

	// validate and finalize
	err = psbt.MaybeFinalizeAll(packet)
	if err != nil {
		return nil, err
	}

	signedTx, err := psbt.Extract(packet)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := signedTx.Serialize(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func signInputP2PK(updater *psbt.Updater, i int, prevPkScript []byte, privKey *btcec.PrivateKey) error {
	// TODO : hashtype always all in p2pk
	hashType := txscript.SigHashAll
	if err := updater.AddInSighashType(hashType, i); err != nil {
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

func signInputP2PKH(updater *psbt.Updater, i int, prevPkScript []byte, privKey *btcec.PrivateKey) error {
	// TODO : hashtype always all in p2pkh
	hashType := txscript.SigHashAll
	if err := updater.AddInSighashType(hashType, i); err != nil {
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

func signInputP2WPKH(updater *psbt.Updater, i int, prevPkScript []byte, amount int64, privKey *btcec.PrivateKey, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	// TODO : hashtype always all in p2wpkh
	hashType := txscript.SigHashAll
	if err := updater.AddInSighashType(hashType, i); err != nil {
		return err
	}

	signature, err := txscript.RawTxInWitnessSignature(updater.Upsbt.UnsignedTx, txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher), i, amount, prevPkScript, hashType, privKey)
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

func signInputP2SH(updater *psbt.Updater, redeemScript []byte, i int, prevPkScript []byte, privKey *btcec.PrivateKey, hashType txscript.SigHashType, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	if err := updater.AddInRedeemScript(redeemScript, i); err != nil {
		return err
	}
	return nil
	// return signInputP2PKH(updater, i, in, redeemScript, privKey, hashType)
}

func signInputMultisig(updater *psbt.Updater, i int, addresses []btcutil.Address, nRequired int, prevPkScript []byte, sigScript, prevScript []byte) error {

	return nil
}

func signInputTaproot(updater *psbt.Updater, privKey *secp256k1.PrivateKey, i int, prevPkScript []byte, prevOutFetcher txscript.PrevOutputFetcher) error {
	var err error

	// key path only
	internalPubKey := schnorr.SerializePubKey(privKey.PubKey())
	updater.Upsbt.Inputs[i].TaprootInternalKey = internalPubKey

	sigHashes := txscript.NewTxSigHashes(updater.Upsbt.UnsignedTx, prevOutFetcher)

	// TODO : hashtype always default in taproot
	hashType := txscript.SigHashDefault
	err = updater.AddInSighashType(hashType, i)
	if err != nil {
		return err
	}
	witness, err := txscript.TaprootWitnessSignature(updater.Upsbt.UnsignedTx, sigHashes,
		i, updater.Upsbt.Inputs[i].WitnessUtxo.Value, prevPkScript, hashType, privKey)
	if err != nil {
		return err
	}
	updater.Upsbt.Inputs[i].TaprootKeySpendSig = witness[0]

	// script path but key path spend
	rootHash := updater.Upsbt.Inputs[i].TaprootMerkleRoot
	if rootHash != nil {
		sig, err := txscript.RawTxInTaprootSignature(updater.Upsbt.UnsignedTx, sigHashes,
			i, updater.Upsbt.Inputs[i].WitnessUtxo.Value, prevPkScript, rootHash, hashType, privKey)
		if err != nil {
			return err
		}
		updater.Upsbt.Inputs[i].TaprootKeySpendSig = sig
	} else {
		if len(updater.Upsbt.Inputs[i].TaprootLeafScript) > 0 {
			// btcd only support one leaf till now
			tapLeaves := updater.Upsbt.Inputs[i].TaprootLeafScript
			taprootScriptSpendSignatures := make([]*psbt.TaprootScriptSpendSig, 0)
			for _, leaf := range tapLeaves {
				tapLeaf := txscript.TapLeaf{
					LeafVersion: leaf.LeafVersion,
					Script:      leaf.Script,
				}
				sig, err := txscript.RawTxInTapscriptSignature(updater.Upsbt.UnsignedTx, sigHashes,
					i, updater.Upsbt.Inputs[i].WitnessUtxo.Value, prevPkScript, tapLeaf, hashType, privKey)
				if err != nil {
					return err
				}
				tapHash := tapLeaf.TapHash()
				tapLeafSignature := &psbt.TaprootScriptSpendSig{
					XOnlyPubKey: internalPubKey,
					LeafHash:    tapHash.CloneBytes(),
					Signature:   sig,
					SigHash:     hashType,
				}
				taprootScriptSpendSignatures = append(taprootScriptSpendSignatures, tapLeafSignature)
			}
			updater.Upsbt.Inputs[i].TaprootInternalKey = nil
			updater.Upsbt.Inputs[i].TaprootKeySpendSig = nil
			// remove duplicate
			updater.Upsbt.Inputs[i].TaprootScriptSpendSig = append(updater.Upsbt.Inputs[i].TaprootScriptSpendSig, taprootScriptSpendSignatures...)
			CheckDuplicateOfUpdater(updater, i)
		}
	}
	return nil
}

func CheckDuplicateOfUpdater(updater *psbt.Updater, index int) {
	signatures := updater.Upsbt.Inputs[index].TaprootScriptSpendSig
	m := map[string]*psbt.TaprootScriptSpendSig{}
	signs := make([]*psbt.TaprootScriptSpendSig, 0)
	for _, v := range signatures {
		key := append(v.XOnlyPubKey, v.LeafHash...)
		keyHex := hex.EncodeToString(key)
		_, ok := m[keyHex]
		if !ok {
			m[keyHex] = v
			signs = append(signs, v)
		}
	}
	updater.Upsbt.Inputs[index].TaprootScriptSpendSig = signs
}

// for verifying the signatures
func ValidateTx(psbtRaw []byte) error {
	// Decode PSBT
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtRaw), false)
	if err != nil {
		return fmt.Errorf("failed to decode PSBT: %v", err)
	}
	prevOutputFetcher := PsbtPrevOutputFetcher(packet)
	for i, input := range packet.Inputs {
		if input.WitnessUtxo == nil && input.NonWitnessUtxo == nil {
			return fmt.Errorf("missing UTXO data for input %d", i)
		}

		// Select UTXO script and value
		var pkScript []byte
		var value int64
		if input.WitnessUtxo != nil {
			pkScript = input.WitnessUtxo.PkScript
			value = input.WitnessUtxo.Value
		} else {
			utxoTx := input.NonWitnessUtxo
			prevOut := utxoTx.TxOut[i]
			pkScript = prevOut.PkScript
			value = prevOut.Value
		}
		// Validate each partial signature
		for _, sig := range input.PartialSigs {
			_ = sig
			// Recreate script engine
			engine, err := txscript.NewEngine(pkScript, packet.UnsignedTx, i, txscript.StandardVerifyFlags, nil, nil, value, prevOutputFetcher)
			if err != nil {
				return fmt.Errorf("failed to create script engine for input %d: %v", i, err)
			}

			// Validate signature
			if err := engine.Execute(); err != nil {
				return fmt.Errorf("signature validation failed for input %d: %v", i, err)
			}
		}
	}

	fmt.Println("All signatures are valid")
	return nil
}
