package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	accent = lipgloss.Color("#89b4fa")

	titleStyle    = lipgloss.NewStyle().Bold(true)
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Bold(true)
	valueStyle    = lipgloss.NewStyle()
	hintStyle     = lipgloss.NewStyle().Faint(true).Italic(true)
	dividerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	itemStyle     = lipgloss.NewStyle()
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
)

func (m model) View() string {
	if m.step == -1 {
		return titleStyle.Render("🧪 Bitcoin Transaction Builder") + "\n\n" +
			errorStyle.Render(m.errorMsg)
	}

	header := titleStyle.Render("🧪 Bitcoin Transaction Builder")
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("────────────────────────────")

	var content string
	if m.step == 0 {
		content = lipgloss.JoinVertical(lipgloss.Left,
			renderStatus(m),
			"",
			labelStyle.Render("SELECT THE NETWORK"),
			renderNetList(m),
			hintStyle.Render("Use ↑/↓ or mouse to move, Enter to select, q to quit"),
			dividerStyle.Render(strings.Repeat("─", 40)),
		)
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left,
			renderStatus(m), "",
			renderInput(m),
		)
	}
	return header + "\n" + divider + "\n" + content
}
func renderNetList(m model) string {
	var b strings.Builder
	idx := m.netList.Index()
	for i, it := range m.netList.Items() {
		name := it.(choiceItem).title
		if i == idx {
			b.WriteString(selectedStyle.Render("➤ · "+name) + "\n")
		} else {
			b.WriteString(itemStyle.Render("  · "+name) + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderStatus(m model) string {
	toMap := formatToMap(m.toList, m.amountList)
	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("STATUS"),
		fmt.Sprintf("🌐 Network:   %s", valueStyle.Render(m.net)),
		fmt.Sprintf("📤 From:      %s", valueStyle.Render(m.from)),
		toMap,
		fmt.Sprintf("🔑 Privatekey: %s", strings.Repeat("X", len(m.privateKey))),
	)
}

func renderInput(m model) string {
	instructions := buildInputInstruction(m.step)

	inputLine := inputStyle.Render(m.inputBuffer + "▌")

	var preview string
	if m.step == 2 && m.inputBuffer != "" {
		preview = fmt.Sprintf("➡️ Preview: %s", m.inputBuffer)
	}

	var errorView string
	if m.errorMsg != "" {
		errorView = errorStyle.Render("⚠️ " + m.errorMsg)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		instructions,
		"",
		inputLine,
		preview,
		lipgloss.NewStyle().Italic(true).Render("Type input and press Enter ⏎"),
		"",
		errorView,
	)
}

func formatToMap(addresses []string, amounts []int64) string {
	if len(addresses) == 0 {
		return "📥 To/Amount: (none)"
	}
	var sb strings.Builder
	sb.WriteString("📥 To/Amount:\n   {\n")
	for i := range addresses {
		sb.WriteString(fmt.Sprintf("     %s: %d,\n", addresses[i], amounts[i]))
	}
	return strings.TrimRight(sb.String(), ",\n") + "\n   }"
}

func buildInputInstruction(step int) string {
	switch step {
	case 0:
		return labelStyle.Render("Enter the Network:") + "\n(btc, btc-testnet3, btc-testnet4)"
	case 1:
		return labelStyle.Render("Enter the 'From' address:")
	case 2:
		return labelStyle.Render("Enter 'To' address and amount:") + "\n(e.g., bc1... 10000) or 'done'"
	case 3:
		return labelStyle.Render("Enter the Fund address (optional):")
	default:
		return ""
	}
}
