package main

import (
	"encoding/hex"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rabbitprincess/btctxbuilder/transaction"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

func (m model) transfer() tea.Msg {
	utxos, err := m.client.GetUTXOWithRawTx(m.from)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to fetch UTXOs: %s", err))
	}
	feeEstimate, err := m.client.FeeEstimate()
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to fetch fee estimate: %s", err))
	}
	fee := max(0.00001, feeEstimate["6"])

	toMap := make(map[string]int64)
	for i := 0; i < len(m.toList); i++ {
		toMap[m.toList[i]] = m.amountList[i]
	}
	params := types.GetParams(types.Network(m.net))

	psbtPacket, err := transaction.NewTransferTx(params, utxos, m.from, toMap, "", fee)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to create transaction: %s", err))
	}

	privKey, err := hex.DecodeString(m.privateKey)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to decode private key: %s", err))
	}
	transaction.SignTx(params, psbtPacket, privKey)

	rawTx, err := types.EncodePsbtToRawTx(psbtPacket)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to encode PSBT to raw transaction: %s", err))
	}

	res, err := m.client.BroadcastTx(utils.HexEncode(rawTx))
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to broadcast transaction: %s", err))
	}

	return resultMsg{txid: res}
}
