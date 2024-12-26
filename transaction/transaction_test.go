package transaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
	"github.com/stretchr/testify/require"
)

// p2pkh
// PrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
// PrivKey, _ := hex.DecodeString(PrivKeyHex)
// PubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
// Address := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"

// p2wpkh
// PrivKeyHex :=  887ad33f247a7df59f1bf61b6aa69ab2a537c0708d7d6fb6614e10511fca377b
// PubKeyHex :=  023c53ee7749c3466415bd8f8b644227b4eb4eaf2339abbb0f1e44e035ea06b21f
// Address :=  tb1q307vt2zz3f66hhs90p0le2pp6r53tvyqnzsy42

// p2tr
// PrivKeyHex := "49b8dbd365939908d920ab74aec8ec9cb3b7d49d252e1aec3ef59bed0f801acc"
// PrivKey, _ := hex.DecodeString(PrivKeyHex)
// Address := "tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal"

// DONE : p2pkh / p2wpkh
// TODO : p2pk / np2wpkh / p2tr
// p2sh, p2wsh np2wsh...??

func TestTransfer(t *testing.T) {
	for _, test := range []struct {
		net         types.Network
		fromPrivKey string
		fromPubKey  string
		fromAddress string
		toAddress   string
		toAmount    int64
	}{
		// p2pkh to p2wpkh
		// {
		// 	types.BTC_Signet,
		// 	"a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409",
		// 	"0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c",
		// 	"n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8",
		// 	"tb1q307vt2zz3f66hhs90p0le2pp6r53tvyqnzsy42",
		// 	1500,
		// },

		// p2wpkh to p2tr
		// {
		// 	types.BTC_Signet,
		// 	"887ad33f247a7df59f1bf61b6aa69ab2a537c0708d7d6fb6614e10511fca377b",
		// 	"023c53ee7749c3466415bd8f8b644227b4eb4eaf2339abbb0f1e44e035ea06b21f",
		// 	"tb1q307vt2zz3f66hhs90p0le2pp6r53tvyqnzsy42",
		// 	"tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal",
		// 	350,
		// },

	} {
		c := client.NewClient(test.net)
		psbtPacket, err := NewTransferTx(c, test.fromAddress, map[string]int64{test.toAddress: test.toAmount}, test.fromAddress)
		require.NoError(t, err)

		fromPrivKey := utils.MustDecode(test.fromPrivKey)
		signedPacket, err := SignTx(c.GetParams(), psbtPacket, fromPrivKey)

		// verify tx
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

}
