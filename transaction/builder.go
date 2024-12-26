package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/types"
)

type BuilderOpt func(*TxBuilder) error

func WithVersion(version int32) BuilderOpt {
	return func(t *TxBuilder) error {
		t.version = version
		return nil
	}
}

func WithFundAddress(address string) BuilderOpt {
	return func(t *TxBuilder) error {
		t.fundAddress = address
		return nil
	}
}

type TxBuilder struct {
	version int32
	client  *client.Client
	params  *chaincfg.Params

	inputs  TxInputs
	outputs TxOutputs

	fundAddress string

	msgTx  *wire.MsgTx
	packet *psbt.Packet
}

func NewTxBuilder(client *client.Client, opts ...BuilderOpt) *TxBuilder {
	builder := &TxBuilder{
		version: wire.TxVersion,
		params:  client.GetParams(),
		client:  client,
	}

	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func (t *TxBuilder) Build() (*psbt.Packet, error) {
	if len(t.inputs) == 0 && len(t.outputs) == 0 {
		return nil, fmt.Errorf("PSBT packet must contain at least one input or output")
	}

	txIns, err := t.inputs.ToWire()
	if err != nil {
		return nil, err
	}
	outputs, err := t.outputs.ToWire()
	if err != nil {
		return nil, err
	}

	t.msgTx = wire.NewMsgTx(t.version)
	for _, in := range txIns {
		t.msgTx.AddTxIn(in)
	}
	for _, out := range outputs {
		t.msgTx.AddTxOut(out)
	}

	if t.fundAddress != "" {
		err = t.FundRawTransaction()
		if err != nil {
			return nil, err
		}
	}

	p, err := psbt.NewFromUnsignedTx(t.msgTx)
	if err != nil {
		return nil, err
	}

	err = t.decorateTxInputs(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (t *TxBuilder) decorateTxInputs(packet *psbt.Packet) error {
	for i := range packet.Inputs {
		txInput := t.inputs[i]

		switch txInput.AddrType {
		case types.P2PK, types.P2PKH:
			addInputInfoNonSegWit(&packet.Inputs[i], txInput)
		case types.P2WPKH, types.P2WPKH_NESTED:
			addInputInfoSegWitV0(&packet.Inputs[i], txInput)
		case types.P2TR:
			addInputInfoSegWitV1(&packet.Inputs[i], txInput)
		default:
			return fmt.Errorf("not support address type %s", txInput.AddrType)
		}

	}
	return nil
}

func addInputInfoNonSegWit(in *psbt.PInput, txInput *TxInput) {
	in.NonWitnessUtxo = txInput.tx

	// Include the derivation path for each input.
	// in.Bip32Derivation = []*psbt.Bip32Derivation{
	// 	derivationInfo,
	// }
}

// addInputInfoSegWitV0 adds the UTXO and BIP32 derivation info for a SegWit v0
// PSBT input (p2wkh, np2wkh) from the given wallet information.
func addInputInfoSegWitV0(in *psbt.PInput, txInput *TxInput) {

	// As a fix for CVE-2020-14199 we have to always include the full
	// non-witness UTXO in the PSBT for segwit v0.
	in.NonWitnessUtxo = txInput.tx

	// To make it more obvious that this is actually a witness output being
	// spent, we also add the same information as the witness UTXO.
	in.WitnessUtxo = txInput.prevVout
	in.SighashType = txscript.SigHashAll

	// Include the derivation path for each input.
	// in.Bip32Derivation = []*psbt.Bip32Derivation{
	// 	derivationInfo,
	// }

	// For nested P2WKH we need to add the redeem script to the input,
	// otherwise an offline wallet won't be able to sign for it. For normal
	// P2WKH this will be nil.
	if txInput.AddrType == types.P2WPKH_NESTED {
		// TODO : test!!
		in.RedeemScript = in.WitnessScript
	}
}

// addInputInfoSegWitV0 adds the UTXO and BIP32 derivation info for a SegWit v1
// PSBT input (p2tr) from the given wallet information.
func addInputInfoSegWitV1(in *psbt.PInput, txInput *TxInput) {

	// For SegWit v1 we only need the witness UTXO information.
	in.WitnessUtxo = txInput.prevVout
	in.SighashType = txscript.SigHashDefault

	// Include the derivation path for each input in addition to the
	// taproot specific info we have below.
	// in.Bip32Derivation = []*psbt.Bip32Derivation{
	// 	derivationInfo,
	// }

	// Include the derivation path for each input.
	// in.TaprootBip32Derivation = []*psbt.TaprootBip32Derivation{{
	// 	XOnlyPubKey:          derivationInfo.PubKey[1:],
	// 	MasterKeyFingerprint: derivationInfo.MasterKeyFingerprint,
	// 	Bip32Path:            derivationInfo.Bip32Path,
	// }}
}
