package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type AddrType string

const (
	P2PKH         AddrType = "p2pkh"        // non segwit
	P2WPKH        AddrType = "p2wpkh"       // native segwit
	P2WPKH_NESTED AddrType = "p2pkh-nested" // nested segwit

	P2SH         AddrType = "p2sh"         // non segwit
	P2WSH        AddrType = "p2wsh"        // native segwit
	P2WSH_NESTED AddrType = "p2wsh-nested" // nested segwit

	TAPROOT AddrType = "taproot" // taproot
)

func PubKeyToAddr(publicKey []byte, addrType AddrType, net Network) (address string, err error) {
	netParams := GetParams(net)
	switch addrType {
	case P2PKH:
		p2pkh, err := btcutil.NewAddressPubKey(publicKey, netParams)
		if err != nil {
			return "", err
		}
		return p2pkh.EncodeAddress(), nil
	case P2WPKH:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		return p2wpkh.EncodeAddress(), nil
	case P2WPKH_NESTED:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		p2sh, err := btcutil.NewAddressScriptHash(redeemScript, netParams)
		if err != nil {
			return "", err
		}
		return p2sh.EncodeAddress(), nil
	case TAPROOT:
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		p2tr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], netParams)
		if err != nil {
			return "", err
		}
		return p2tr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)
	}
}

func ScriptToAddr(script []byte, addrType AddrType, net Network) (address string, err error) {
	netParams := GetParams(net)
	switch addrType {
	case P2SH:
		// OP_HASH160 <ScriptHash> OP_EQUAL
		if len(script) != 23 || script[0] != txscript.OP_HASH160 {
			return "", fmt.Errorf("invalid P2SH script")
		}
		p2sh, err := btcutil.NewAddressScriptHashFromHash(script[2:22], netParams)
		if err != nil {
			return "", err
		}
		return p2sh.EncodeAddress(), nil
	case P2WSH:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return "", fmt.Errorf("invalid native segwit script")
		}
		witnessProgram := script[2:]
		p2wsh, err := btcutil.NewAddressWitnessScriptHash(witnessProgram, netParams)
		if err != nil {
			return "", err
		}
		return p2wsh.EncodeAddress(), nil
	case P2WSH_NESTED:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return "", fmt.Errorf("invalid nested segwit script")
		}
		redeemScript, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_0).
			AddData(script[2:]).
			Script()
		if err != nil {
			return "", fmt.Errorf("failed to create redeem script: %w", err)
		}
		redeemScriptHash := btcutil.Hash160(redeemScript)
		p2shAddr, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, netParams)
		if err != nil {
			return "", fmt.Errorf("failed to create SegWit Nested address: %w", err)
		}
		return p2shAddr.EncodeAddress(), nil
	case TAPROOT:
		// OP_1 <32-byte-TweakHash>
		if len(script) != 34 || script[0] != txscript.OP_1 {
			return "", fmt.Errorf("invalid Taproot script")
		}
		taprootKey := script[2:]
		p2tr, err := btcutil.NewAddressTaproot(taprootKey, netParams)
		if err != nil {
			return "", err
		}
		return p2tr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)
	}
}
