package types

import (
	"bytes"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

type TxType string

const (
	Transfer TxType = "transfer"
	FeeBump  TxType = "feebump"

	Script TxType = "script"
)

func DecodePSBT(psbtStr string) (*psbt.Packet, error) {
	var err error
	var psbtRaw []byte

	isHex := utils.IsHex(psbtStr)
	if isHex {
		psbtRaw, err = utils.Decode(psbtStr)
		if err != nil {
			return nil, err
		}
	} else {
		psbtRaw = []byte(psbtStr)
	}
	p, err := psbt.NewFromRawBytes(bytes.NewReader(psbtRaw), !isHex)
	if err != nil {
		return nil, err
	}
	return p, nil
}
