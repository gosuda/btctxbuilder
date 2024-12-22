package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubKeyToAddr(t *testing.T) {
	network := BTC_Testnet3
	params := GetParams(network)
	pubKeyHex := "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f"
	publicKey, err := hex.DecodeString(pubKeyHex)
	assert.NoError(t, err)

	p2pk, err := PubKeyToAddr(publicKey, P2PK, params)
	require.NoError(t, err)
	// base58 encoded compressed public key
	assert.Equal(t, "zbRapgvpp4xSYvt8oeuzBc9qfZh2UfAgQ6r218xhCQxe", p2pk)

	p2pkh, err := PubKeyToAddr(publicKey, P2PKH, params)
	require.NoError(t, err)
	assert.Equal(t, "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", p2pkh)

	p2wpkh, err := PubKeyToAddr(publicKey, P2WPKH, params)
	require.NoError(t, err)
	assert.Equal(t, "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", p2wpkh)

	nestedP2wpkh, err := PubKeyToAddr(publicKey, P2WPKH_NESTED, params)
	require.NoError(t, err)
	assert.Equal(t, "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", nestedP2wpkh)

	p2tr, err := PubKeyToAddr(publicKey, TAPROOT, params)
	require.NoError(t, err)
	assert.Equal(t, "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", p2tr)
}

func TestAddrType(t *testing.T) {
	network := BTC_Testnet3
	params := GetParams(network)
	addr, err := DecodeAddress("mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", params)
	require.NoError(t, err)
	p2pkh, err := GetAddressType(addr)
	require.NoError(t, err)
	require.Equal(t, P2PKH, p2pkh)

	addr, err = DecodeAddress("tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", params)
	require.NoError(t, err)
	p2wpkh, err := GetAddressType(addr)
	require.NoError(t, err)
	require.Equal(t, P2WPKH, p2wpkh)

	addr, err = DecodeAddress("2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", params)
	require.NoError(t, err)
	nestedP2wpkh, err := GetAddressType(addr)
	require.NoError(t, err)
	// p2wpkh-nested = p2sh
	require.Equal(t, P2SH, nestedP2wpkh)

	addr, err = DecodeAddress("tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", params)
	require.NoError(t, err)
	p2tr, err := GetAddressType(addr)
	require.NoError(t, err)
	require.Equal(t, TAPROOT, p2tr)

}

func TestGenerateAddress(t *testing.T) {
	network := BTC_Testnet3
	params := GetParams(network)

	priv, err := secp256k1.GeneratePrivateKey()
	require.NoError(t, err)
	privKey := priv.Serialize()

	pub := priv.PubKey()
	pubKey := pub.SerializeCompressed()
	addressP2PKH, err := PubKeyToAddr(pubKey, P2PKH, params)
	require.NoError(t, err)

	fmt.Println("Private Key: ", hex.EncodeToString(privKey))
	fmt.Println("Public Key: ", hex.EncodeToString(pubKey))
	fmt.Println("P2PKH Address: ", addressP2PKH)

}
