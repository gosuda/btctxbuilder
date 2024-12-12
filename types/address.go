package types

type AddrType string

const (
	P2PK          AddrType = "p2pk"          // non segwit
	P2PKH         AddrType = "p2pkh"         // non segwit
	P2WPKH        AddrType = "p2wpkh"        // native segwit
	P2WPKH_NESTED AddrType = "p2wpkh-nested" // nested segwit

	P2SH         AddrType = "p2sh"         // non segwit
	P2WSH        AddrType = "p2wsh"        // native segwit
	P2WSH_NESTED AddrType = "p2wsh-nested" // nested segwit

	TAPROOT AddrType = "taproot" // taproot

	Invalid AddrType = ""
)
