package types

type TxType string

const (
	Transfer TxType = "transfer"
	FeeBump  TxType = "feebump"

	Script TxType = "script"
)
