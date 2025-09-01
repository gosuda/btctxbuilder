package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gosuda/btctxbuilder/address"
	"github.com/gosuda/btctxbuilder/types"

	"github.com/gosuda/btctxbuilder/transaction"
)

/* ---------- newAddress flow ---------- */

func (m model) generateNewAddress(typ string) (model, tea.Cmd) {
	var at types.AddrType
	switch strings.ToLower(typ) {
	case string(types.P2PK):
		at = types.P2PK
	case string(types.P2PKH):
		at = types.P2PKH
	case string(types.P2WPKH):
		at = types.P2WPKH
	case string(types.P2WPKH_NESTED):
		at = types.P2WPKH_NESTED
	case string(types.P2TR):
		at = types.P2TR
	default:
		return m, returnError("Unsupported address type: " + typ)
	}

	privHex, pubHex, addr, err := address.GenerateAddress(at)
	if err != nil {
		return m, returnError(fmt.Sprintf("GenerateAddress failed: %v", err))
	}
	m.resultAddr = addr
	m.resultPubHex = fmt.Sprintf("%s", pubHex)
	m.resultPrivHex = fmt.Sprintf("%s", privHex)
	m.bannerKind = "success"
	m.banner = fmt.Sprintf("New %s address generated.\n⚠️  WARNING: This data is NOT stored by the app. Please copy and store it in a secure place!", at)
	m.step = 11
	return m, nil
}

func (m model) transfer() tea.Msg {
	toMap := make(map[string]int64)
	for i := 0; i < len(m.toList); i++ {
		toMap[m.toList[i]] = m.amountList[i]
	}

	signer, err := types.NewECDSASigner(m.privateKey)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to decode private key: %s", err))
	}

	txid, err := transaction.BroadcastTx(m.client, m.from, toMap, signer.Sign, signer.PubKey())
	if err != nil {
		return errorMsg(err.Error())
	}
	return resultMsg{txid: txid}
}
