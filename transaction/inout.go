package transaction

type Vin struct {
	txId         string
	vOut         uint32
	redeemScript string
	address      string
	amount       int64
}

type Vout struct {
	address string
	script  string
	amount  int64
}
