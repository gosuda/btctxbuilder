package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type AddrType string

const (
	P2PK          AddrType = "p2pk"          // non segwit
	P2PKH         AddrType = "p2pkh"         // non segwit
	P2WPKH        AddrType = "p2wpkh"        // native segwit
	P2WPKH_NESTED AddrType = "p2wpkh-nested" // nested segwit

	P2SH         AddrType = "p2sh"         // non segwit
	P2WSH        AddrType = "p2wsh"        // native segwit
	P2WSH_NESTED AddrType = "p2wsh-nested" // nested segwit

	TAPROOT AddrType = "taproot" // taproot

	Invalid AddrType = ""
)

func PubKeyToAddr(publicKey []byte, addrType AddrType, params *chaincfg.Params) (address string, err error) {
	switch addrType {
	case P2PK:
		addr, err := btcutil.NewAddressPubKey(publicKey, params)
		if err != nil {
			return "", err
		}
		return base58.Encode(addr.ScriptAddress()), nil
	case P2PKH:
		addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey), params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case P2WPKH:
		address, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), params)
		if err != nil {
			return "", err
		}
		return address.EncodeAddress(), nil
	case P2WPKH_NESTED:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), params)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		addr, err := btcutil.NewAddressScriptHash(redeemScript, params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case TAPROOT:
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		addr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)

	}
}

func ScriptToAddr(script []byte, addrType AddrType, params *chaincfg.Params) (address string, err error) {
	switch addrType {
	case P2SH:
		// OP_HASH160 <ScriptHash> OP_EQUAL
		if len(script) != 23 || script[0] != txscript.OP_HASH160 {
			return "", fmt.Errorf("invalid P2SH script")
		}
		addr, err := btcutil.NewAddressScriptHashFromHash(script[2:22], params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case P2WSH:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return "", fmt.Errorf("invalid native segwit script")
		}
		witnessProgram := script[2:]
		addr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram, params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
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
		addr, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case TAPROOT:
		// OP_1 <32-byte-TweakHash>
		if len(script) != 34 || script[0] != txscript.OP_1 {
			return "", fmt.Errorf("invalid Taproot script")
		}
		taprootKey := script[2:]
		addr, err := btcutil.NewAddressTaproot(taprootKey, params)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)
	}
}

func DecodeAddress(address string, params *chaincfg.Params) (addr btcutil.Address, err error) {
	return btcutil.DecodeAddress(address, params)
}

func GetAddressType(addr btcutil.Address) (addrType AddrType, err error) {

	switch addr.(type) {
	case *btcutil.AddressPubKey:
		return P2PK, nil
	case *btcutil.AddressPubKeyHash:
		return P2PKH, nil
	case *btcutil.AddressScriptHash:
		return P2SH, nil
	case *btcutil.AddressWitnessPubKeyHash:
		return P2WPKH, nil
	case *btcutil.AddressWitnessScriptHash:
		return P2WSH, nil
	case *btcutil.AddressTaproot:
		return TAPROOT, nil
	default:
		return Invalid, fmt.Errorf("unsupported address type: %T", addr)
	}
}
