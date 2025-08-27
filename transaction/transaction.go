package transaction

import (
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/gosuda/btctxbuilder/types"
)

// fee is expected in sat/vB. signer may be nil to return an unsigned PSBT.
func NewTransferTx(
	params *chaincfg.Params,
	utxos []*types.Utxo,
	fromAddress string,
	toAddress map[string]int64,
	signer types.Signer,
	pubkey []byte,
	fee float64,
) (*psbt.Packet, error) {

	builder := NewTxBuilder(params).
		FeeRate(fee).
		From(fromAddress).
		Change(fromAddress).
		ToMap(toAddress).
		SelectInputs(utxos).
		Build().
		SignWith(signer, pubkey)

	// check err
	if err := builder.Err(); err != nil {
		return nil, err
	}
	return builder.Packet(), nil
}

func NewRunestoneEdictTx(params *chaincfg.Params, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) {

}
