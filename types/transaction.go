package types

type TransactionType string

const (
	Transfer TransactionType = "transfer"
	Multisig TransactionType = "multisig"
	Timelock TransactionType = "timelock"

	Ordinals TransactionType = "ordinals"
)
