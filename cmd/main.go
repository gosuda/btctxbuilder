package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rabbitprincess/btctxbuilder/client"
	"github.com/rabbitprincess/btctxbuilder/transaction"
	"github.com/rabbitprincess/btctxbuilder/types"
	"github.com/rabbitprincess/btctxbuilder/utils"
)

type model struct {
	step        int
	net         string
	from        string
	toList      []string
	amountList  []int64
	fundAddress string
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
		switch msg.Type {
		case tea.KeyEnter:
			input := strings.TrimSpace(m.inputBuffer)
			m.inputBuffer = ""
			switch m.step {
			case 0:
				client, err := client.NewClient(types.Network(input))
				if err != nil {
					return m, func() tea.Msg {
						return errorMsg(fmt.Sprintf("Failed to create client: %s", err))
					}
				}
				m.net = input
				m.client = client
				m.step++
			case 1:
				if input == "" {
					return m, func() tea.Msg {
						return errorMsg("From address cannot be empty.")
					}
				}
				m.from = input
				m.step++
			case 2:
				if input == "done" {
					if len(m.toList) == 0 {
						return m, func() tea.Msg {
							return errorMsg("At least one recipient is required.")
						}
					}
					m.step++
				} else {
					var address string
					var amount int64
					_, err := fmt.Sscanf(input, "%s %d", &address, &amount)
					if err != nil || address == "" || amount <= 0 {
						return m, func() tea.Msg {
							return errorMsg("Invalid format. Use: [address] [amount]")
						}
					}
					m.toList = append(m.toList, address)
					m.amountList = append(m.amountList, amount)
				}
			case 3:
				if input == "" {
					m.fundAddress = m.from
				} else {
					m.fundAddress = input
				}
				return m, m.transfer
			}
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		default:
			m.inputBuffer += string(msg.Runes)
		}
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

func (m model) View() string {
	header := titleStyle.Render("🧪 Bitcoin Transaction Builder")
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("────────────────────────────")

	// status block
	toMap := "📥 To/Amount:\n   {\n"
	for i := range m.toList {
		toMap += fmt.Sprintf("     %s: %d,\n", m.toList[i], m.amountList[i])
	}
	if len(m.toList) > 0 {
		toMap = strings.TrimRight(toMap, ",\n") + "\n   }"
	} else {
		toMap = "📥 To/Amount: (none)"
	}

	status := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("🌐 Network:     %s", m.net),
		fmt.Sprintf("📤 From:        %s", m.from),
		toMap,
		fmt.Sprintf("💰 Fund Addr:   %s", m.fundAddress),
	)

	statusBox := lipgloss.NewStyle().
		Width(90).
		Border(lipgloss.RoundedBorder()).
		Background(lipgloss.Color("234")).
		Padding(1, 1).
		Render(status)

	// input section
	var instructions string
	switch m.step {
	case 0:
		instructions = highlightStyle.Render("Enter the Network:") + "\n(btc, btc-testnet3, btc-testnet4)"
	case 1:
		instructions = highlightStyle.Render("Enter the 'From' address:")
	case 2:
		instructions = highlightStyle.Render("Enter 'To' address and amount:") + "\n(e.g., bc1... 10000) or 'done'"
	case 3:
		instructions = highlightStyle.Render("Enter the Fund address (optional):")
	case -1:
		return lipgloss.NewStyle().Padding(1, 2).Render(header + "\n\n" + titleStyle.Render(m.errorMsg))
	}

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 1).
		Width(80).
		Render(m.inputBuffer + "▌")

	inputHint := lipgloss.NewStyle().
		Italic(true).
		Render("Type input and press Enter ⏎")

	var preview string
	if m.step == 2 && m.inputBuffer != "" {
		preview = inputStyle.Render(fmt.Sprintf("➡️ Preview: %s", m.inputBuffer))
	}

	var errorView string
	if m.errorMsg != "" {
		errorView = errorStyle.Render("⚠️ " + m.errorMsg)
	}

	inputs := lipgloss.JoinVertical(lipgloss.Left,
		instructions,
		"",
		inputBox,
		preview,
		inputHint,
		"",
		errorView,
	)

	inputBoxWrap := lipgloss.NewStyle().
		Width(90).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 1).
		Render(inputs)

	// outer wrap box
	content := lipgloss.JoinVertical(lipgloss.Left,
		statusBox,
		inputBoxWrap,
	)

	mainBox := lipgloss.NewStyle().
		Render(content)

	return lipgloss.NewStyle().
		Render(header + "\n" + divider + "\n\n" + mainBox)
}

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

	psbtPacket, err := transaction.NewTransferTx(m.client.GetParams(), utxos, m.from, toMap, "", fee)
	if err != nil {
		return errorMsg(fmt.Sprintf("Failed to create transaction: %s", err))
	}

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

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	inputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("13")).Foreground(lipgloss.Color("0")).Bold(true)
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
