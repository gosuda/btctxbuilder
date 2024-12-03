package types

import "github.com/btcsuite/btcd/txscript"

type ScriptType string

const (
	ScriptP2PK     ScriptType = "p2pk"
	ScriptP2PKH    ScriptType = "p2pkh"
	ScriptP2SH     ScriptType = "p2sh"
	ScriptP2WPKH   ScriptType = "p2wpkh"
	ScriptP2WSH    ScriptType = "p2wsh"
	ScriptP2TR     ScriptType = "p2tr"
	ScriptMultisig ScriptType = "multisig"
	ScriptOpReturn ScriptType = "opreturn"

	ScriptUnknown ScriptType = "unknown"
)

func ParseScriptType(pkScript []byte) ScriptType {
	switch {
	case txscript.IsPayToPubKey(pkScript): // P2PK
		return ScriptP2PK
	case txscript.IsPayToPubKeyHash(pkScript): // P2PKH
		return ScriptP2PKH
	case txscript.IsPayToScriptHash(pkScript): // P2SH
		return ScriptP2SH
	case txscript.IsPayToWitnessPubKeyHash(pkScript): // P2WPKH
		return ScriptP2WPKH
	case txscript.IsPayToWitnessScriptHash(pkScript): // P2WSH
		return ScriptP2WSH
	case txscript.IsPayToTaproot(pkScript): // P2TR
		return ScriptP2TR
	case IsMultiSigScript(pkScript): // Multisig
		return ScriptMultisig
	case txscript.IsNullData(pkScript): // OP_RETURN
		return ScriptOpReturn
	default:
		return ScriptUnknown
	}
}

func IsMultiSigScript(script []byte) bool {
	ok, _ := txscript.IsMultisigScript(script)
	return ok
}
