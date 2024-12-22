package transaction

import (
	"bytes"
	"testing"

	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignPsbtTx(t *testing.T) {
	networkType := types.BTC_Testnet3

	c := client.NewClient(networkType)
	txBuild := NewTxBuilder(c)
	txBuild.inputs.AddInputTransfer("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 0)
	txBuild.outputs.AddOutputTransfer(c.Params, "mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.outputs.AddOutputTransfer(c.Params, "mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)
	pubKeyMap := make(map[int]string)
	pubKeyMap[0] = "022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933"
	packet, err := txBuild.Build()
	require.Nil(t, err)
	// signatureMap := make(map[int]string)
	// for i, h := range hashes {
	// 	privateBytes, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	// 	require.Nil(t, err)
	// 	prvKey, _ := btcec.PrivKeyFromBytes(privateBytes)
	// 	sign := ecdsa.Sign(prvKey, RemoveZeroHex(h))
	// 	signatureMap[i] = hex.EncodeToString(sign.Serialize())
	// }
	var buf bytes.Buffer
	err = packet.Serialize(&buf)
	require.NoError(t, err)
	psbtRaw := buf.Bytes()

	rawTx, err := SignTx(networkType, psbtRaw, utils.MustDecode("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	require.Nil(t, err)
	assert.Equal(t, "01000000012258bcdf8563608ffdf86aca3580b0f4830d5b92d3a1cf906c32e6c2f5232c0b010000006a47304402206bdac667fb3d6f1a62e0b0d1123a5caa58d8c0fd95c2a2c8cd091374960a871702204f301e6883866570ce309573e569d6a32a44386af5bf928b5f9e1dcd7e2dd0ed0121022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933ffffffff0208cf0000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac10270000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac00000000", rawTx)
}
