package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type AddrType string

const (
	P2PKH         AddrType = "p2pkh"
	P2SH          AddrType = "p2sh"
	SEGWIT_NATIVE AddrType = "segwit_native"
	SEGWIT_NESTED AddrType = "segwit_nested"
	TAPROOT       AddrType = "taproot"
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
	case SEGWIT_NATIVE:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		return p2wpkh.EncodeAddress(), nil
	case SEGWIT_NESTED:
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
	case P2PKH:
		// OP_DUP OP_HASH160 <PubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
		if len(script) != 25 || script[0] != txscript.OP_DUP || script[1] != txscript.OP_HASH160 {
			return "", fmt.Errorf("invalid P2PKH script")
		}
		p2pkh, err := btcutil.NewAddressPubKeyHash(script[3:23], netParams)
		if err != nil {
			return "", err
		}
		return p2pkh.EncodeAddress(), nil
	case P2SH:
		if len(script) != 23 || script[0] != txscript.OP_HASH160 {
			return "", fmt.Errorf("invalid P2SH script")
		}
		p2sh, err := btcutil.NewAddressScriptHashFromHash(script[2:22], netParams)
		if err != nil {
			return "", err
		}
		return p2sh.EncodeAddress(), nil
	case SEGWIT_NATIVE:
		if len(script) < 2 || script[0] != txscript.OP_0 {
			return "", fmt.Errorf("invalid native segwit script")
		}
		witnessProgram := script[2:]
		p2wsh, err := btcutil.NewAddressWitnessScriptHash(witnessProgram, netParams)
		if err != nil {
			return "", err
		}
		return p2wsh.EncodeAddress(), nil
	case TAPROOT:
		if len(script) < 34 || script[0] != txscript.OP_1 {
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
