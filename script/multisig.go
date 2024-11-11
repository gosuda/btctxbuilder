package script

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/rabbitprincess/btctxbuilder/types"
)

func EncodeMultiSigScript(network types.Network, pubKeys [][]byte, nRequired int) ([]byte, error) {
	if nRequired <= 1 {
		nRequired = len(pubKeys) - 1
	}

	addrPubKeys := make([]*btcutil.AddressPubKey, 0, len(pubKeys))
	for _, pubKey := range pubKeys {
		addrPubkey, err := btcutil.NewAddressPubKey(pubKey, types.GetParams(network))
		if err != nil {
			return nil, err
		}
		addrPubKeys = append(addrPubKeys, addrPubkey)
	}

	script, err := txscript.MultiSigScript(addrPubKeys, nRequired)
	if err != nil {
		return nil, err
	}
	return script, nil
}

func DecodeMultiSigScript(script []byte) ([][]byte, error) {
	// nRequired, pubKeys, _, err := txscript.MakeScriptTokenizer(script)
	// if err != nil {
	// return 0, nil, err
	// }
	// return nRequired, pubKeys, nil
	return nil, nil // TODO
}
