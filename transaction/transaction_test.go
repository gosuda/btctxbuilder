package transaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
	"github.com/stretchr/testify/require"
)

// p2pkh
// fromPrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
// fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
// fromPubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
// fromAddress := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"

// p2wpkh

// p2tr
// fromPrivKeyHex := "49b8dbd365939908d920ab74aec8ec9cb3b7d49d252e1aec3ef59bed0f801acc"
// fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
// fromAddress := "tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal"

func TestTransferP2PKH(t *testing.T) {
	fromPrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
	fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
	fromPubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
	fromAddress := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"
	toAddress := "tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal"
	var toAmount int64 = 1000

	net := types.BTC_Signet
	c := client.NewClient(net)
	psbtPacket, err := NewTransferTx(c, fromAddress, map[string]int64{toAddress: toAmount}, fromAddress)
	require.NoError(t, err)

	signedPacket, err := SignTx(c.GetParams(), psbtPacket, fromPrivKey)

	// verify tx
	pub, err := secp256k1.ParsePubKey(utils.MustDecode(fromPubKey))
	require.NoError(t, err)

	valid, err := VerifyTx(c.GetParams(), signedPacket, pub)
	require.NoError(t, err)
	fmt.Println(valid)

	require.NoError(t, err)
	signedTxRaw, err := types.EncodePsbtToRawTx(signedPacket)
	require.NoError(t, err)
	signedTxHex := hex.EncodeToString(signedTxRaw)
	fmt.Println(signedTxHex)

	// txid, err := btcclient.BroadcastTx(signedTxHex)
	// require.NoError(t, err)
	// fmt.Println("txid:", txid)

	newTx, err := client.DecodeRawTransaction(signedTxHex)
	require.NoError(t, err)

	jsonNewTx, _ := json.MarshalIndent(newTx, "", "\t")
	fmt.Println(string(jsonNewTx))
}
