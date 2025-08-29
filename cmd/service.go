package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gosuda/btctxbuilder/types"

	"github.com/gosuda/btctxbuilder/transaction"
)

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
