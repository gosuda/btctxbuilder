package transaction

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"

	"github.com/gosuda/btctxbuilder/types"
)

// -----------------------------------------------------------------------------
// context
// -----------------------------------------------------------------------------

type txctx struct {
	params     *chaincfg.Params
	feeRate    float64
	fromAddr   string
	changeAddr string
	utxos      []*types.Utxo
	Inputs     TxInputs
	Outputs    TxOutputs
	pkt        *psbt.Packet

	errs []error // accumulate all errors
}

func (c *txctx) addErr(err error) {
	if err != nil {
		c.errs = append(c.errs, err)
	}
}

func (c *txctx) Err() error {
	if len(c.errs) == 0 {
		return nil
	}
	return errors.Join(c.errs...)
}

// -----------------------------------------------------------------------------
// states
// -----------------------------------------------------------------------------

func NewTxBuilder(params *chaincfg.Params) BInit { return BInit{&txctx{params: params}} }

type BInit struct{ c *txctx }   // only initial configuration
type BDraft struct{ c *txctx }  // add inputs/outputs/select
type BBuilt struct{ c *txctx }  // psbt built (unsigned)
type BSigned struct{ c *txctx } // psbt signed

// error getters (available on all states)
func (b BInit) Err() error   { return b.c.Err() }
func (b BDraft) Err() error  { return b.c.Err() }
func (b BBuilt) Err() error  { return b.c.Err() }
func (b BSigned) Err() error { return b.c.Err() }

// -----------------------------------------------------------------------------
// BInit: configuration (only here!)
// -----------------------------------------------------------------------------

func (b BInit) From(addr string) BInit        { b.c.fromAddr = addr; return b }
func (b BInit) Change(addr string) BInit      { b.c.changeAddr = addr; return b }
func (b BInit) FeeRate(feeRate float64) BInit { b.c.feeRate = feeRate; return b }

// move into draft by adding first thing (input or output)
func (b BInit) AddInput(raw *wire.MsgTx, vout uint32, amt int64, addr string) BDraft {
	b.c.addErr(b.c.Inputs.AddInput(b.c.params, raw, vout, amt, addr))
	return BDraft{b.c}
}
func (b BInit) To(addr string, amt int64) BDraft {
	b.c.addErr(b.c.Outputs.AddOutputTransfer(b.c.params, addr, amt))
	return BDraft{b.c}
}

// -----------------------------------------------------------------------------
// BDraft: keep adding inputs/outputs or auto-select inputs
// -----------------------------------------------------------------------------

func (b BDraft) AddInput(raw *wire.MsgTx, vout uint32, amt int64, addr string) BDraft {
	b.c.addErr(b.c.Inputs.AddInput(b.c.params, raw, vout, amt, addr))
	return b
}
func (b BDraft) To(addr string, amt int64) BDraft {
	b.c.addErr(b.c.Outputs.AddOutputTransfer(b.c.params, addr, amt))
	return b
}

// SelectInputs picks UTXOs to cover current Outputs (fee finalized in funding step).
func (b BDraft) SelectInputs(utxos []*types.Utxo) BDraft {
	c := b.c
	need := int64(c.Outputs.AmountTotal())

	sel, rest, err := SelectUtxo(utxos, need)
	if err != nil {
		c.addErr(err)
		return b
	}

	for _, u := range sel {
		c.addErr(c.Inputs.AddInput(c.params, u.RawTx, u.Vout, u.Value, c.fromAddr))
	}
	c.utxos = rest
	if c.changeAddr == "" {
		c.changeAddr = c.fromAddr
	}
	return b
}

// Build: single error boundary
func (b BDraft) Build() (BBuilt, error) {
	if err := b.c.Err(); err != nil {
		return BBuilt{}, err
	}

	msg := wire.NewMsgTx(wire.TxVersion)

	ins, err := b.c.Inputs.ToWire()
	if err != nil {
		return BBuilt{}, err
	}
	for _, in := range ins {
		msg.AddTxIn(in)
	}

	outs, err := b.c.Outputs.ToWire()
	if err != nil {
		return BBuilt{}, err
	}
	for _, out := range outs {
		msg.AddTxOut(out)
	}

	// finalize fee + add change (may also add more inputs via inputs pointer)
	if b.c.changeAddr != "" {
		if err := FundRawTransaction(
			b.c.params,
			msg,
			&b.c.Inputs,
			b.c.Outputs,
			b.c.changeAddr,
			b.c.feeRate,
			b.c.utxos,
			b.c.fromAddr,
			SelectUtxo,
		); err != nil {
			return BBuilt{}, err
		}
	}

	pkt, err := psbt.NewFromUnsignedTx(msg)
	if err != nil {
		return BBuilt{}, err
	}

	if err := DecorateTxInputs(pkt, b.c.Inputs); err != nil {
		return BBuilt{}, err
	}

	b.c.pkt = pkt
	return BBuilt{b.c}, nil
}

// -----------------------------------------------------------------------------
// BBuilt: unsigned packet helpers
// -----------------------------------------------------------------------------

func (b BBuilt) Packet() *psbt.Packet { return b.c.pkt }

func (b BBuilt) UnsignedRawTx() ([]byte, error) {
	if err := b.c.Err(); err != nil {
		return nil, err
	}
	if b.c.pkt == nil || b.c.pkt.UnsignedTx == nil {
		return nil, fmt.Errorf("no unsigned tx: call Build() first")
	}
	var buf bytes.Buffer
	if err := b.c.pkt.UnsignedTx.Serialize(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b BBuilt) SignWith(sign types.Signer, pubkey []byte) (BSigned, error) {
	if b.c.pkt == nil {
		return BSigned{}, fmt.Errorf("packet is nil: call Build() first")
	}
	pkt, err := SignTx(b.c.params, b.c.pkt, sign, pubkey)
	if err != nil {
		return BSigned{}, err
	}
	b.c.pkt = pkt
	return BSigned{b.c}, nil
}

// -----------------------------------------------------------------------------
// BSigned: final raw tx (broadcastable)
// -----------------------------------------------------------------------------

func (b BSigned) Packet() *psbt.Packet { return b.c.pkt }

func (b BSigned) RawTx() ([]byte, error) {
	if b.c.pkt == nil {
		return nil, fmt.Errorf("packet is nil")
	}
	finalTx, err := psbt.Extract(b.c.pkt)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := finalTx.Serialize(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
