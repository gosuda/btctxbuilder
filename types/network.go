package types

import (
	"log"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
)

type Network string

const (
	BTC               Network = "btc"
	BTC_Testnet3      Network = "btc-testnet3"
	BTC_Regressionnet Network = "btc-regtest"
	BTC_Signet        Network = "btc-signet"

	DGB  Network = "dgb"  // digibyte
	QTUM Network = "qtum" // qtum
	RVN  Network = "rvn"  // raven
	BTG  Network = "btg"  // bitcoin gold
	BCH  Network = "bch"  // bitcoin cash
	DOGE Network = "doge" // dogecoin
)

var (
	netParams = map[Network]*chaincfg.Params{}
)

func init() {
	netParams = map[Network]*chaincfg.Params{
		BTC:               getBTCMainNetParams(),
		BTC_Testnet3:      getBTCTestNetParams(),
		BTC_Regressionnet: getBTCRegresstionNetParams(),
		BTC_Signet:        getBTCSignetParams(),

		DGB:  getDGBMainNetParams(),
		QTUM: getQTUMMainNetParams(),
		RVN:  getRVNMainNetParams(),
		BTG:  getBTGMainNetParams(),
		BCH:  getBCHmainNetParams(),
		DOGE: getDOGEMainNetParams(),
	}
}

func GetParams(net Network) *chaincfg.Params {
	if param, ok := netParams[net]; ok {
		return param
	}
	log.Fatalf("network not supported [%s]", net)
	return nil
}

// getBTCMainNetParams BTC
func getBTCMainNetParams() *chaincfg.Params {
	return &chaincfg.MainNetParams
}

func getBTCTestNetParams() *chaincfg.Params {
	return &chaincfg.TestNet3Params
}

func getBTCRegresstionNetParams() *chaincfg.Params {
	return &chaincfg.RegressionNetParams
}

func getBTCSignetParams() *chaincfg.Params {
	return &chaincfg.SigNetParams
}

// getDGBMainNetParams DGB
func getDGBMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0xdab6c3fa

	// Address encoding magics
	params.PubKeyHashAddrID = 30 // base58 prefix: D
	params.ScriptHashAddrID = 63 // base58 prefix: 3
	params.Bech32HRPSegwit = "dgb"
	return &params
}

// GetQTUMMainNetParams QTUM
func getQTUMMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0xf1cfa6d3

	// Address encoding magics
	params.PubKeyHashAddrID = 58 // base58 prefix: Q
	params.ScriptHashAddrID = 50 // base58 prefix: P
	params.Bech32HRPSegwit = "qc"

	return &params
}

// getRVNMainNetParams RVN
func getRVNMainNetParams() *chaincfg.Params {
	params := chaincfg.MainNetParams
	params.Net = 0x4e564152

	// Address encoding magics
	params.PubKeyHashAddrID = 60  // base58 prefix: R
	params.ScriptHashAddrID = 122 // base58 prefix: r
	return &params
}

// getBTGMainNetParams BTG
func getBTGMainNetParams() *chaincfg.Params {
	mainnetparams := chaincfg.MainNetParams
	mainnetparams.Net = 0x446d47e1

	// Address encoding magics
	mainnetparams.PubKeyHashAddrID = 38 // base58 prefix: G
	mainnetparams.ScriptHashAddrID = 23 // base58 prefix: A

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	// see https://github.com/satoshilabs/slips/blob/master/slip-0173.md
	mainnetparams.Bech32HRPSegwit = "btg"

	return &mainnetparams
}

// getBCHmainNetParams BCH
func getBCHmainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xe8f3e1e3

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 0
	mainNetParams.ScriptHashAddrID = 5
	return &mainNetParams
}

// getLTCMainNetParams LTC
func getLTCMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xdbb6c0fb

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 48
	mainNetParams.ScriptHashAddrID = 50
	mainNetParams.Bech32HRPSegwit = "ltc"
	return &mainNetParams
}

// getDASHMainNetParams DASH
func getDASHMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xbd6b0cbf

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 76
	mainNetParams.ScriptHashAddrID = 16
	return &mainNetParams
}

// getDOGEMainNetParams DOGE
func getDOGEMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = 0xc0c0c0c0

	// Address encoding magics
	mainNetParams.PubKeyHashAddrID = 30
	mainNetParams.ScriptHashAddrID = 22 // base58 prefix: 9
	return &mainNetParams
}

func NewZECAddr(pubBytes []byte) string {
	version := []byte{0x1C, 0xB8}
	return NewOldAddr(version, btcutil.Hash160(pubBytes))
}

func NewOldAddr(version []byte, data []byte) string {
	var buf []byte
	buf = append(buf, version[1:]...)
	buf = append(buf, data...)
	return base58.CheckEncode(buf, version[0])
}

const (
	zecNet = 0x6427e924
)

// getZECMainNetParams ZEC
func getZECMainNetParams() *chaincfg.Params {
	mainNetParams := chaincfg.MainNetParams
	mainNetParams.Net = zecNet

	mainNetParams.PubKeyHashAddrID = 0x1C
	mainNetParams.ScriptHashAddrID = 0xBD

	return &mainNetParams
}
