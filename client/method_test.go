package client

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/stretchr/testify/require"
)

func TestGetBestBlock(t *testing.T) {
	client := NewClient(types.BTC)
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
		fmt.Println(tx.Txid)
		// rawTx, err := client.GetRawTx(tx.Txid)
		require.NoError(t, err)
		// fmt.Println("rawTx:", string(rawTx))
		for _, vout := range tx.Vout {
			fmt.Printf("%s ", vout.ScriptpubkeyType)
		}
		fmt.Println()
	}
}

// transaction example
// p2pkh : ef796f3cef041768d37a34a469d72e5c91de568f963eae6daf3480fe8405e2ed
// v0_p2wpkh : 6c9f507a64cfec9ef96de41680af40c84607d71b62eac7f7f2a406a597c8c582
// p2sh : 6216b12925f9bf817679e4cbaae35e1f5b8da997dc8b12603c6de7dd965af5c1
// v0_p2wsh : ca31304e07751c96dfc9c48812a3404759fb31c89694efc27cbe1a72d1d439d8
// v1_p2tr : dcf80b086238982841bfc382a5a567c8f6898878db44d9da0d3726edc7bb7211
func TestGetTx(t *testing.T) {
	client := NewClient(types.BTC)
	tx, err := client.GetTx("6c9f507a64cfec9ef96de41680af40c84607d71b62eac7f7f2a406a597c8c582")
	require.NoError(t, err)

	txJson, _ := json.MarshalIndent(tx, "", "\t")
	fmt.Println(string(txJson))

	for _, vin := range tx.Vin {
		if vin.Prevout != nil {
			fmt.Println("prev vout value :", vin.Prevout.Value)
			fmt.Println("prev vout scriptpubkey :", vin.Prevout.Scriptpubkey)
			fmt.Println("prev vout scriptpubkey asm :", vin.Prevout.ScriptpubkeyAsm)
			fmt.Println("type :", vin.Prevout.ScriptpubkeyType)
		}

		fmt.Println("vin script sig :", vin.Scriptsig)
		fmt.Println("vin script sig asm :", vin.ScriptsigAsm)
		fmt.Println("vin witness :", vin.Witness)
		fmt.Println("vin sequence :", vin.Sequence)
		fmt.Println()
	}

	for _, vout := range tx.Vout {
		fmt.Println("vout scriptpubkey :", vout.Scriptpubkey)
		fmt.Println("vout scriptpubkey asm :", vout.ScriptpubkeyAsm)
		fmt.Println("vout scriptpubkey type :", vout.ScriptpubkeyType)
		fmt.Println("vout scriptpubkey address :", vout.ScriptpubkeyAddress)
		fmt.Println("vout value :", vout.Value)
		fmt.Println()
	}

}

// fromPrivKeyHex := "a6018c89646f3c7596516544602283135e8d6e5b31421e335b91b86ae9c76409"
// fromPrivKey, _ := hex.DecodeString(fromPrivKeyHex)
// fromPubKey := "0248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925c"
// fromAddress := "n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8"
// toAddress := "tb1plt7057su6z39qjqtnvnnw7d6htdwulqm93mtpddj5wcetwxcv2nsm6geal"
func TestGetBalance(t *testing.T) {
	client := NewClient(types.BTC_Signet)
	addr, err := client.GetAddress("n368zCWREFiRRX7icJRBb6n8nMsjJjNVK8")
	require.NoError(t, err)
	fmt.Println(addr.Address)
	fmt.Println("funded sat :", addr.ChainStats.FundedTxoSum)
	fmt.Println("spent sat :", addr.ChainStats.SpentTxoSum)
	fmt.Println("balance :", addr.ChainStats.FundedTxoSum-addr.ChainStats.SpentTxoSum)
	fmt.Println("tx count :", addr.ChainStats.TxCount)

}
func TestGetUtxo(t *testing.T) {
	client := NewClient(types.BTC)
	utxos, err := client.GetUTXO("bc1qwzrryqr3ja8w7hnja2spmkgfdcgvqwp5swz4af4ngsjecfz0w0pqud7k38")
	require.NoError(t, err)

	fmt.Println(utxos)

	for _, utxo := range utxos {
		fmt.Println(utxo)
	}
}

func TestFeeEstimate(t *testing.T) {
	client := NewClient(types.BTC)
	fee, err := client.FeeEstimate()
	require.NoError(t, err)

	fmt.Println(fee)
}

func TestBroadCastTx(t *testing.T) {
	client := NewClient(types.BTC_Signet)
	tx := "0100000001d9736e55bb7a1ba989455c7181e88afde41cb190364e858db430d401121bca13010000006a47304402205e1ba81844e9b5a7e37e89f1bdbd7abbcd07e19b6d39f4a6dcb100fff22b81bf02204ccbd2eca2f33a3777e39a22e5a19aa7bdb382e570c4e58719332fcd74ab176f01210248d7c76f23e387bb151e6094590eb8f7777a8efbea9d0a5ddd1ea1833fa3925cffffffff02dc050000000000001600148bfcc5a8428a75abde05785ffca821d0e915b080ac120000000000001976a914eca14b26ef6056bf1011137061a5ffdbecba4c6188ac00000000"

	rawTx, err := DecodeRawTransaction(tx)
	require.NoError(t, err)
	fmt.Println("txid :", rawTx.TxID())
	for _, txIn := range rawTx.TxIn {
		fmt.Println("\tvin hash  :", txIn.PreviousOutPoint.Hash)
		fmt.Println("\tvin index :", txIn.PreviousOutPoint.Index)
		fmt.Println("\tvin sig :", txIn.SignatureScript)
		fmt.Println("\tvin witness :", txIn.Witness)
		fmt.Println("\tvin sequence :", txIn.Sequence)
	}
	for _, txOut := range rawTx.TxOut {
		fmt.Println("\tvout script :", txOut.PkScript)
		fmt.Println("\tvout value  :", txOut.Value)
	}

	res, err := client.BroadcastTx(tx)
	require.NoError(t, err)
	fmt.Println("result:", res)
}
