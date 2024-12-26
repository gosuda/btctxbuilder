package types

import (
	"bytes"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
)

func EncodePsbt(packet *psbt.Packet) ([]byte, error) {
	var buf bytes.Buffer
	err := packet.Serialize(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodePsbt(rawPacket []byte) (*psbt.Packet, error) {
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(rawPacket), false)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

func EncodePsbtToRawTx(packet *psbt.Packet) ([]byte, error) {
	signedTx, err := psbt.Extract(packet)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := signedTx.Serialize(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeRawTx(rawTx []byte) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	err := tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		return nil, err
	}
	return tx, nil
}
