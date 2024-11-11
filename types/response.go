package types

type Block struct {
	ID                string  `json:"id"`
	Height            uint64  `json:"height"`
	Version           int     `json:"version"`
	Timestamp         int     `json:"timestamp"`
	TxCount           int     `json:"tx_count"`
	Size              int     `json:"size"`
	Weight            int     `json:"weight"`
	MerkleRoot        string  `json:"merkle_root"`
	Previousblockhash string  `json:"previousblockhash"`
	Mediantime        int     `json:"mediantime"`
	Nonce             int     `json:"nonce"`
	Bits              int     `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
}

type Transaction struct {
	Txid     string      `json:"txid"`
	Version  int         `json:"version"`
	Locktime int         `json:"locktime"`
	Vin      []Vin       `json:"vin"`
	Vout     []Vout      `json:"vout"`
	Size     int         `json:"size"`
	Weight   int         `json:"weight"`
	Fee      int         `json:"fee"`
	Status   BlockStatus `json:"status"`
}

type Vin struct {
	Txid         string   `json:"txid"`
	Vout         uint32   `json:"vout"`
	Prevout      any      `json:"prevout"`
	Scriptsig    string   `json:"scriptsig"`
	ScriptsigAsm string   `json:"scriptsig_asm"`
	Witness      []string `json:"witness"`
	IsCoinbase   bool     `json:"is_coinbase"`
	Sequence     int64    `json:"sequence"`

	// calculated fields
	Amount  int64  `json:"-"`
	Address string `json:"-"`
}

type Vout struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
	ScriptpubkeyType    string `json:"scriptpubkey_type"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address,omitempty"`
	Value               int64  `json:"value"`

	// calculated fields
	Amount    int64  `json:"-"`
	Address   string `json:"-"`
	PublicKey string `json:"-"`
}

type Utxo struct {
	Txid   string      `json:"txid"`
	Vout   uint32      `json:"vout"`
	Status BlockStatus `json:"status"`
	Value  int64       `json:"value"`
}

type BlockStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int    `json:"block_time"`
}

type FeeEstimate map[string]float64
