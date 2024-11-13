package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type AddrType string

const (
	P2PKH         AddrType = "p2pkh"         // non segwit
	P2WPKH        AddrType = "p2wpkh"        // native segwit
	P2WPKH_NESTED AddrType = "p2wpkh-nested" // nested segwit

	P2SH         AddrType = "p2sh"         // non segwit
	P2WSH        AddrType = "p2wsh"        // native segwit
	P2WSH_NESTED AddrType = "p2wsh-nested" // nested segwit

	TAPROOT AddrType = "taproot" // taproot
)

func PubKeyToAddr(publicKey []byte, addrType AddrType, net Network) (address btcutil.Address, err error) {
	netParams := GetParams(net)

	switch addrType {
	case P2PKH:
		address, err = btcutil.NewAddressPubKey(publicKey, netParams)
	case P2WPKH:
		address, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
	case P2WPKH_NESTED:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return nil, err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return nil, err
		}
		address, err = btcutil.NewAddressScriptHash(redeemScript, netParams)
	case TAPROOT:
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return nil, err
		}
		address, err = btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], netParams)
	default:
		err = fmt.Errorf("address type not supported | %s", addrType)
	}
	return address, err
}

func ScriptToAddr(script []byte, addrType AddrType, net Network) (address btcutil.Address, err error) {
	netParams := GetParams(net)
	switch addrType {
	case P2SH:
		// OP_HASH160 <ScriptHash> OP_EQUAL
		if len(script) != 23 || script[0] != txscript.OP_HASH160 {
			return nil, fmt.Errorf("invalid P2SH script")
		}
		address, err = btcutil.NewAddressScriptHashFromHash(script[2:22], netParams)
	case P2WSH:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return nil, fmt.Errorf("invalid native segwit script")
		}
		witnessProgram := script[2:]
		address, err = btcutil.NewAddressWitnessScriptHash(witnessProgram, netParams)
	case P2WSH_NESTED:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return nil, fmt.Errorf("invalid nested segwit script")
		}
		redeemScript, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_0).
			AddData(script[2:]).
			Script()
		if err != nil {
			return nil, fmt.Errorf("failed to create redeem script: %w", err)
		}
		redeemScriptHash := btcutil.Hash160(redeemScript)
		address, err = btcutil.NewAddressScriptHashFromHash(redeemScriptHash, netParams)
	case TAPROOT:
		// OP_1 <32-byte-TweakHash>
		if len(script) != 34 || script[0] != txscript.OP_1 {
			return nil, fmt.Errorf("invalid Taproot script")
		}
		taprootKey := script[2:]
		address, err = btcutil.NewAddressTaproot(taprootKey, netParams)
	default:
		err = fmt.Errorf("address type not supported | %s", addrType)
	}
	return address, err
}
