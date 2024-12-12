package address

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/txscript"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func PubKeyToAddr(publicKey []byte, addrType types.AddrType, net types.Network) (address string, err error) {
	netParams := types.GetParams(net)

	switch addrType {
	case types.P2PK:
		addr, err := btcutil.NewAddressPubKey(publicKey, netParams)
		if err != nil {
			return "", err
		}
		return base58.Encode(addr.ScriptAddress()), nil
	case types.P2PKH:
		addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case types.P2WPKH:
		address, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		return address.EncodeAddress(), nil
	case types.P2WPKH_NESTED:
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), netParams)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		addr, err := btcutil.NewAddressScriptHash(redeemScript, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case types.TAPROOT:
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		addr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)

	}
}

func ScriptToAddr(script []byte, addrType types.AddrType, net types.Network) (address string, err error) {
	netParams := types.GetParams(net)
	switch addrType {
	case types.P2SH:
		// OP_HASH160 <ScriptHash> OP_EQUAL
		if len(script) != 23 || script[0] != txscript.OP_HASH160 {
			return "", fmt.Errorf("invalid P2SH script")
		}
		addr, err := btcutil.NewAddressScriptHashFromHash(script[2:22], netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case types.P2WSH:
		// OP_0 <32-byte-ScriptHash>
		if len(script) != 34 || script[0] != txscript.OP_0 {
			return "", fmt.Errorf("invalid native segwit script")
		}
		witnessProgram := script[2:]
		addr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case types.P2WSH_NESTED:
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
		addr, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case types.TAPROOT:
		// OP_1 <32-byte-TweakHash>
		if len(script) != 34 || script[0] != txscript.OP_1 {
			return "", fmt.Errorf("invalid Taproot script")
		}
		taprootKey := script[2:]
		addr, err := btcutil.NewAddressTaproot(taprootKey, netParams)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("address type not supported | %s", addrType)
	}
}

func DecodeAddress(address string, net types.Network) (addr btcutil.Address, err error) {
	netParams := types.GetParams(net)
	return btcutil.DecodeAddress(address, netParams)
}

func GetAddressType(address string, net types.Network) (addrType types.AddrType, err error) {
	addr, err := DecodeAddress(address, net)
	if err != nil {
		return types.Invalid, err
	}

	switch addr.(type) {
	case *btcutil.AddressPubKey:
		return types.P2PK, nil
	case *btcutil.AddressPubKeyHash:
		return types.P2PKH, nil
	case *btcutil.AddressScriptHash:
		return types.P2SH, nil
	case *btcutil.AddressWitnessPubKeyHash:
		return types.P2WPKH, nil
	case *btcutil.AddressWitnessScriptHash:
		return types.P2WSH, nil
	case *btcutil.AddressTaproot:
		return types.TAPROOT, nil
	default:
		return types.Invalid, fmt.Errorf("unsupported address type: %T", addr)
	}
}
