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
