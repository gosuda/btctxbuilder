package transaction

import (
	"fmt"
	"sort"

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
	fundAddress string,
	signer types.Signer, // nil = return unsigned PSBT
	pubkey []byte, // required if signer != nil
	fee float64, // sat/vB
) (*psbt.Packet, error) {
	if len(toAddress) == 0 {
		return nil, fmt.Errorf("no outputs: toAddress is empty")
	}
	if fundAddress == "" {
		fundAddress = fromAddress
	}

	// Init configuration (only allowed on BInit)
	init := NewTxBuilder(params).
		From(fromAddress).
		Change(fundAddress).
		FeeRate(fee)

	// Deterministic output order
	addrs := make([]string, 0, len(toAddress))
	for a := range toAddress {
		addrs = append(addrs, a)
	}
	sort.Strings(addrs)

	// Move into draft by adding outputs
	var draft BDraft
	for i, a := range addrs {
		amt := toAddress[a]
		if i == 0 {
			draft = init.To(a, amt)
		} else {
			draft = draft.To(a, amt)
		}
	}

	// Select inputs from the provided UTXO pool (fee finalized in Build via FundRawTransaction)
	draft = draft.SelectInputs(utxos)

	// Build unsigned PSBT
	built, err := draft.Build()
	if err != nil {
		return nil, err
	}

	// Optionally sign
	if signer == nil {
		return built.Packet(), nil
	}
	signed, err := built.SignWith(signer, pubkey)
	if err != nil {
		return nil, err
	}
	return signed.Packet(), nil
}

func NewRunestoneEdictTx(params *chaincfg.Params, utxos []*types.Utxo, fromAddress string, toAddress map[string]int64, fundAddress string) {

}
