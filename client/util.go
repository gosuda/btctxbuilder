package client

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

func DecodeRawTransaction(rawTx string) (*wire.MsgTx, error) {
	// Parse the raw transaction
	var rawTxBytes []byte
	if utils.IsHex(rawTx) {
		rawTxBytes = utils.MustDecode(rawTx)
	} else {
		rawTxBytes = []byte(rawTx)
	}

	msgTx := wire.NewMsgTx(wire.TxVersion)
	err := msgTx.Deserialize(bytes.NewReader(rawTxBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	return msgTx, nil
}
