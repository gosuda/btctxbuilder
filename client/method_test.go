package client

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

func TestGetBestBlock(t *testing.T) {
	client := NewClient(Mainnet)
	height, err := client.BestBlockHeight()
	require.NoError(t, err)
	fmt.Println("height:", height)

	hash, err := client.BestBlockHash()
	require.NoError(t, err)

	fmt.Println("hash:", hash)

	block, err := client.GetBlock(hash)
	require.NoError(t, err)

	fmt.Println("block:", block)

	txs, err := client.GetBlockTx(hash, 0)
	require.NoError(t, err)
	for _, tx := range txs {
		fmt.Println(tx)
		rawTx, err := client.GetRawTx(tx.Txid)
		require.NoError(t, err)
		fmt.Println("rawTx:", string(rawTx))
	}
}

func TestGetUtxo(t *testing.T) {
	client := NewClient(Mainnet)
	utxos, err := client.GetUTXO("bc1qwzrryqr3ja8w7hnja2spmkgfdcgvqwp5swz4af4ngsjecfz0w0pqud7k38")
	require.NoError(t, err)

	fmt.Println(utxos)

	for _, utxo := range utxos {
		fmt.Println(utxo)
	}
}

func TestBroadCastTx(t *testing.T) {
	client := NewClient(Mainnet)
	raw, err := hex.DecodeString("0100000000010181318803dc1a178fce37d628cf832e8bb18e94492cf109caa232c40f9e68c2f20100000000ffffffff02404c1f0400000000160014973d7c4a508283a3727aa0c512594a24bfd99824d2f17d0500000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d040047304402207c49f592a903ba568afe0b58ce76a76853a5af993907c1c20b11b38c20f4a566022042105e7baf2565c59d119dd763b70e30f720943960d6e6eb21ca16afdefc18f80147304402201ebb6849245a4b8e67c9ac411f613442b7f5515ad0027e91fa36e341462aa2ac02201d0b04731ea7e47142ec86605c1d3a67528a347ab22d41a62e0371a9555e17e4016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000")
	require.NoError(t, err)
	fmt.Println(string(raw))

	tx := "0100000000010181318803dc1a178fce37d628cf832e8bb18e94492cf109caa232c40f9e68c2f20100000000ffffffff02404c1f0400000000160014973d7c4a508283a3727aa0c512594a24bfd99824d2f17d0500000000220020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d040047304402207c49f592a903ba568afe0b58ce76a76853a5af993907c1c20b11b38c20f4a566022042105e7baf2565c59d119dd763b70e30f720943960d6e6eb21ca16afdefc18f80147304402201ebb6849245a4b8e67c9ac411f613442b7f5515ad0027e91fa36e341462aa2ac02201d0b04731ea7e47142ec86605c1d3a67528a347ab22d41a62e0371a9555e17e4016952210375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c2103a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff2103c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f88053ae00000000"
	rawTx, err := DecodeRawTransaction([]byte(tx))
	require.NoError(t, err)
	fmt.Println(rawTx)

	hash, err := client.BroadcastTx(tx)
	require.NoError(t, err)
	fmt.Println("hash:", hash)
}

func DecodeRawTransaction(rawTxBytes []byte) (*wire.MsgTx, error) {

	// Parse the raw transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)
	err := msgTx.Deserialize(hex.NewDecoder(bytes.NewReader(rawTxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	return msgTx, nil
}
