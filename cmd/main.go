package main

// import (
// 	"fmt"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/rabbitprincess/btctxbuilder/client"
// 	"github.com/rabbitprincess/btctxbuilder/transaction"
// 	"github.com/rabbitprincess/btctxbuilder/types"
// 	"github.com/rabbitprincess/btctxbuilder/utils"
// )

// type model struct {
// 	step        int
// 	net         string
// 	from        string
// 	toAddress   map[string]int64
// 	fundAddress string
// 	client      *client.Client
// 	errorMsg    string
// 	inputBuffer string
// }

// // Message types
// type inputMsg string
// type errorMsg string

// type resultMsg struct {
// 	psbt string
// }

// // Initial model
// func initialModel() model {
// 	return model{
// 		step:      0,
// 		toAddress: make(map[string]int64),
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// // Update function
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyEnter:
// 			// 엔터를 누르면 입력값 처리
// 			input := m.inputBuffer
// 			m.inputBuffer = "" // 입력 버퍼 초기화
// 			switch m.step {
// 			case 0:
// 				if input != "btc" && input != "btc-testnet3" && input != "btc-signet" {
// 					return m, func() tea.Msg { return errorMsg("Invalid network. Use: btc, btc-testnet3, or btc-signet.") }
// 				}
// 				client := client.NewClient(types.Network(input))
// 				m.net = input
// 				m.client = client
// 				m.step++
// 			case 1:
// 				if input == "" {
// 					return m, func() tea.Msg { return errorMsg("'From' address cannot be empty.") }
// 				}
// 				m.from = input
// 				m.step++
// 			case 2:
// 				var address string
// 				var amount int64
// 				_, err := fmt.Sscanf(input, "%s %d", &address, &amount)
// 				if err != nil || address == "" || amount <= 0 {
// 					return m, func() tea.Msg { return errorMsg("Invalid address format. Use: address amount (amount > 0).") }
// 				}
// 				m.toAddress[address] = amount
// 				m.step++
// 			case 3:
// 				if input == "" {
// 					m.fundAddress = m.from
// 				} else {
// 					m.fundAddress = input
// 				}
// 				return m, m.transfer
// 			}
// 		case tea.KeyBackspace:
// 			if len(m.inputBuffer) > 0 {
// 				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
// 			}
// 		default:
// 			m.inputBuffer += string(msg.Runes)
// 		}
// 	case errorMsg:
// 		m.errorMsg = string(msg)
// 		return m, tea.Batch()
// 	case resultMsg:
// 		m.step = -1
// 		m.errorMsg = fmt.Sprintf("Transaction created successfully! PSBT: %s", msg.psbt)
// 		return m, nil
// 	}

// 	return m, nil
// }

// var (
// 	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
// 	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
// 	inputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
// 	errorStyleBold = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
// 	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("13")).Foreground(lipgloss.Color("0"))
// )

// func (m model) View() string {
// 	if m.errorMsg != "" {
// 		return fmt.Sprintf("Error: %s\nPress any key to restart.", m.errorMsg)
// 	}

// 	view := ""
// 	switch m.step {
// 	case 0:
// 		view = titleStyle.Render("Enter the Network for the client (btc, btc-testnet3, btc-signet):")
// 	case 1:
// 		view = titleStyle.Render("Enter the 'from' address:")
// 	case 2:
// 		view = inputStyle.Render("Enter 'to' address and amount (e.g., address 100000):")
// 	case 3:
// 		view = inputStyle.Render("Enter the fund address (or press Enter to use 'from' address):")
// 	case -1:
// 		view = errorStyleBold.Render(m.errorMsg)
// 	default:
// 		view = errorStyleBold.Render("An unknown error occurred.")
// 	}

// 	return fmt.Sprintf("%s\n\nInput: %s", view, m.inputBuffer)
// }

// // transfer triggers the NewTransferTx function
// func (m model) transfer() tea.Msg {
// 	utxos, err := m.client.GetUTXO(m.from)
// 	if err != nil {
// 		return errorMsg(fmt.Sprintf("Failed to fetch UTXOs: %s", err))
// 	}

// 	psbtPacket, err := transaction.NewTransferTx(m.client, utxos, m.from, m.toAddress, m.fundAddress)
// 	if err != nil {
// 		return errorMsg(fmt.Sprintf("Failed to create transaction: %s", err))
// 	}

// 	rawTx, err := types.EncodePsbtToRawTx(psbtPacket)
// 	if err != nil {
// 		return errorMsg(fmt.Sprintf("Failed to encode PSBT to raw transaction: %s", err))
// 	}

// 	return resultMsg{psbt: utils.Encode(rawTx)}
// }

// func main() {
// 	// Create a new program
// 	p := tea.NewProgram(initialModel())

// 	// Run the program and handle any errors
// 	if _, err := p.Run(); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		os.Exit(1)
// 	}
// }
