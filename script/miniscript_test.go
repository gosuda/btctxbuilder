package script

import (
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/require"
)

func TestPkNode(t *testing.T) {
	for _, test := range []struct {
		compressed bool
	}{
		{true},
		{false},
	} {
		privKey, err := btcec.NewPrivateKey()
		require.NoError(t, err)
		var pubKey []byte
		if test.compressed {
			pubKey = privKey.PubKey().SerializeCompressed()
		} else {
			pubKey = privKey.PubKey().SerializeUncompressed()
		}

		pkNode := &PkNode{PubKey: pubKey}

		// create script from public key
		builder := txscript.NewScriptBuilder()
		pkNode.ToScript(builder)
		script, err := builder.Script()
		require.NoError(t, err)

		// recover public key from script
		pkNodeNew := &PkNode{}
		err = pkNodeNew.FromScript(script)
		require.NoError(t, err)

		require.Equal(t, pubKey, pkNodeNew.PubKey)
	}
}

func TestPkhNode(t *testing.T) {
	for _, test := range []struct {
		compressed bool
	}{
		{true},
		{false},
	} {
		privKey, err := btcec.NewPrivateKey()
		require.NoError(t, err)
		var pubKeyHash []byte
		if test.compressed {
			pubKey := privKey.PubKey().SerializeCompressed()
			pubKeyHash = btcutil.Hash160(pubKey)
		} else {
			pubKey := privKey.PubKey().SerializeUncompressed()
			pubKeyHash = btcutil.Hash160(pubKey)
		}
		pkhNode := &PkhNode{Hash: pubKeyHash}

		// create script from public key
		builder := txscript.NewScriptBuilder()
		pkhNode.ToScript(builder)
		script, err := builder.Script()
		require.NoError(t, err)

		// recover public key from script
		pkhNodeNew := &PkhNode{}
		err = pkhNodeNew.FromScript(script)
		require.NoError(t, err)

		require.Equal(t, pubKeyHash, pkhNodeNew.Hash)
	}

}
