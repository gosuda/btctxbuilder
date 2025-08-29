package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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

	netList list.Model
}

type choiceItem struct{ title, desc string }

func (i choiceItem) Title() string       { return i.title }
func (i choiceItem) Description() string { return i.desc }
func (i choiceItem) FilterValue() string { return i.title }

type simpleDelegate struct{}

func (d simpleDelegate) Height() int                               { return 1 }
func (d simpleDelegate) Spacing() int                              { return 0 }
func (d simpleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d simpleDelegate) Render(w io.Writer, m list.Model, idx int, it list.Item) {
	item, ok := it.(choiceItem)
	if !ok {
		return
	}
	prefix := "  "
	if idx == m.Index() {
		prefix = "➤ "
	}
	fmt.Fprintf(w, "%s%s\n", prefix, item.title)
}

type errorMsg string
type resultMsg struct{ txid string }

func initialModel() model {
	items := []list.Item{
		choiceItem{"btc", "Bitcoin mainnet"},
		choiceItem{"btc-testnet3", "Legacy testnet3"},
		choiceItem{"btc-testnet4", "New testnet4"},
		choiceItem{"btc-signet", "Signet"},
	}
	l := list.New(items, simpleDelegate{}, 24, 8)
	l.Title = ""
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.FilterInput.Focus()

	return model{step: 0, netList: l}
}

func (m model) Init() tea.Cmd { return nil }

/* ---------- update ---------- */

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "esc":
			return m, tea.Quit
		}

		if m.step == 0 {
			return m.updateNetworkSelect(x)
		}
		return m.updateKeyInput(x)

	case errorMsg:
		m.errorMsg = string(x)
		return m, nil

	case resultMsg:
		m.step = -1
		m.errorMsg = fmt.Sprintf("Transaction successful! txid: %s", x.txid)
		return m, nil
	}

	if m.step == 0 {
		var cmd tea.Cmd
		m.netList, cmd = m.netList.Update(msg)
		return m, cmd
	}
	return m, nil
}

/* ---------- step==0: network select ---------- */

func (m model) updateNetworkSelect(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.Type {
	case tea.KeyEnter:
		if it := m.netList.SelectedItem(); it != nil {
			sel := it.(choiceItem).title
			m.netList.FilterInput.Blur()
			return m.setNetwork(sel)
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.netList, cmd = m.netList.Update(k)
		return m, cmd
	}
}

/* ---------- step>=1: text input ---------- */

func (m model) updateKeyInput(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.Type {
	case tea.KeyEnter:
		input := strings.TrimSpace(m.inputBuffer)
		m.inputBuffer = ""
		return m.handleStep(input)
	case tea.KeyBackspace:
		if n := len(m.inputBuffer); n > 0 {
			m.inputBuffer = m.inputBuffer[:n-1]
		}
	default:
		if len(k.Runes) > 0 {
			m.inputBuffer += string(k.Runes)
		}
	}
	return m, nil
}

func (m model) handleStep(input string) (model, tea.Cmd) {
	switch m.step {
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

/* ---------- business actions ---------- */

func (m model) setNetwork(input string) (model, tea.Cmd) {
	client, err := client.NewClient(types.Network(input))
	if err != nil {
		return m, returnError(fmt.Sprintf("Failed to create client: %s", err))
	}
	m.net = input
	m.client = client
	m.step = 1
	return m, nil
}

func (m model) setFromAddress(input string) (model, tea.Cmd) {
	if input == "" {
		return m, returnError("From address cannot be empty.")
	}
	m.from = input
	m.step = 2
	return m, nil
}

func (m model) addRecipient(input string) (model, tea.Cmd) {
	if input == "done" {
		if len(m.toList) == 0 {
			return m, returnError("At least one recipient is required.")
		}
		m.step = 3
		return m, nil
	}
	var addr string
	var amt int64
	if _, err := fmt.Sscanf(input, "%s %d", &addr, &amt); err != nil || addr == "" || amt <= 0 {
		return m, returnError("Invalid format. Use: [address] [amount]")
	}
	m.toList = append(m.toList, addr)
	m.amountList = append(m.amountList, amt)
	return m, nil
}

func (m *model) setPrivateKey(input string) { m.privateKey = input }

func returnError(msg string) tea.Cmd { return func() tea.Msg { return errorMsg(msg) } }
