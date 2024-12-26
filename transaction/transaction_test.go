package transaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/stretchr/testify/require"
)

// p2pkh
// fromPrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
// fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
// fromPubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
// fromAddress := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"

// p2wpkh
// fromPrivKeyHex :=  887ad33f247a7df59f1bf61b6aa69ab2a537c0708d7d6fb6614e10511fca377b
// fromPubKeyHex :=  023c53ee7749c3466415bd8f8b644227b4eb4eaf2339abbb0f1e44e035ea06b21f
// fromAddress :=  tb1q307vt2zz3f66hhs90p0le2pp6r53tvyqnzsy42

// p2tr
// fromPrivKeyHex := "49b8dbd365939908d920ab74aec8ec9cb3b7d49d252e1aec3ef59bed0f801acc"
// fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
// fromAddress := "tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal"

func TestTransferP2PKH(t *testing.T) {
	fromPrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
	fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
	fromPubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
	fromAddress := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"
	toAddress := "tb1q307vt2zz3f66hhs90p0le2pp6r53tvyqnzsy42"
	var toAmount int64 = 1500

	net := types.BTC_Signet
	c := client.NewClient(net)
	psbtPacket, err := NewTransferTx(c, fromAddress, map[string]int64{toAddress: toAmount}, fromAddress)
	require.NoError(t, err)

	signedPacket, err := SignTx(c.GetParams(), psbtPacket, fromPrivKey)

	// verify tx
	_ = fromPubKey
	// pub, err := secp256k1.ParsePubKey(utils.MustDecode(fromPubKey))
	// require.NoError(t, err)
	// valid, err := VerifyTx(c.GetParams(), signedPacket, pub)
	// require.NoError(t, err)
	// fmt.Println(valid)

	require.NoError(t, err)
	signedTxRaw, err := types.EncodePsbtToRawTx(signedPacket)
	require.NoError(t, err)
	signedTxHex := hex.EncodeToString(signedTxRaw)
	fmt.Println(signedTxHex)

	txid, err := c.BroadcastTx(signedTxHex)
	require.NoError(t, err)
	fmt.Println("txid:", txid)

	newTx, err := client.DecodeRawTransaction(signedTxHex)
	require.NoError(t, err)

	jsonNewTx, _ := json.MarshalIndent(newTx, "", "\t")
	fmt.Println(string(jsonNewTx))
}
