package client

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

func DecodeRawTransaction(rawTxBytes []byte) (*wire.MsgTx, error) {
	// Parse the raw transaction
	if utils.IsHex(string(rawTxBytes)) {
		rawTxBytes = utils.MustDecode(string(rawTxBytes))
	}

	msgTx := wire.NewMsgTx(wire.TxVersion)
	err := msgTx.Deserialize(bytes.NewReader(rawTxBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	return msgTx, nil
}
