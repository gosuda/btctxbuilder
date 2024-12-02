package types

type ScriptType string

const (
	Key          ScriptType = "key"
	RedeemScript ScriptType = "redeemscript"
	MAST         ScriptType = "mast"

	OpReturn ScriptType = "opreturn" // For embedding arbitrary data
	HashLock ScriptType = "hashlock" // Hash lock condition (e.g., Atomic Swap)
)

type TransactionType string

const (
	Transfer TransactionType = "transfer"
	Multisig TransactionType = "multisig"
	Timelock TransactionType = "timelock"

	Ordinals TransactionType = "ordinals"
)
