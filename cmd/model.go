package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gosuda/btctxbuilder/client"
	"github.com/gosuda/btctxbuilder/types"
)

type model struct {
	step        int
	net         string
	from        string
	toList      []string
	amountList  []int64
	privateKey  string
	client      *client.Client
	errorMsg    string
	inputBuffer string
}

type errorMsg string
type resultMsg struct{ txid string }

func initialModel() model {
	return model{step: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.updateKey(msg)
	case errorMsg:
		m.errorMsg = string(msg)
		return m, nil
	case resultMsg:
		m.step = -1
		m.errorMsg = fmt.Sprintf("Transaction successful! txid: %s", msg.txid)
		return m, nil
	}
	return m, nil
}

func (m model) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		input := strings.TrimSpace(m.inputBuffer)
		m.inputBuffer = ""
		return m.handleStep(input)
	case tea.KeyBackspace:
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
	default:
		m.inputBuffer += string(msg.Runes)
	}
	return m, nil
}

func (m model) handleStep(input string) (model, tea.Cmd) {
	switch m.step {
	case 0:
		return m.setNetwork(input)
	case 1:
		return m.setFromAddress(input)
	case 2:
		return m.addRecipient(input)
	case 3:
		m.setPrivateKey(input)
		return m, m.transfer
	}
	return m, nil
}

func (m model) setNetwork(input string) (model, tea.Cmd) {
	client, err := client.NewClient(types.Network(input))
	if err != nil {
		return m, returnError(fmt.Sprintf("Failed to create client: %s", err))
	}
	m.net = input
	m.client = client
	m.step++
	return m, nil
}

func (m model) setFromAddress(input string) (model, tea.Cmd) {
	if input == "" {
		return m, returnError("From address cannot be empty.")
	}
	m.from = input
	m.step++
	return m, nil
}

func (m model) addRecipient(input string) (model, tea.Cmd) {
	if input == "done" {
		if len(m.toList) == 0 {
			return m, returnError("At least one recipient is required.")
		}
		m.step++
		return m, nil
	}
	var address string
	var amount int64
	_, err := fmt.Sscanf(input, "%s %d", &address, &amount)
	if err != nil || address == "" || amount <= 0 {
		return m, returnError("Invalid format. Use: [address] [amount]")
	}
	m.toList = append(m.toList, address)
	m.amountList = append(m.amountList, amount)
	return m, nil
}

func (m *model) setPrivateKey(input string) {
	m.privateKey = input
}

func returnError(msg string) tea.Cmd {
	return func() tea.Msg { return errorMsg(msg) }
}
