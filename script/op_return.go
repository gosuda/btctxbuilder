package script

import "github.com/btcsuite/btcd/txscript"

func OpReturnScript(data []byte) ([]byte, error) {
	return txscript.NullDataScript(data)
}
