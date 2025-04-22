package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	inputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	highlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("13")).Foreground(lipgloss.Color("0")).Bold(true)
)

func (m model) View() string {
	if m.step == -1 {
		return lipgloss.NewStyle().Padding(1, 2).Render(
			titleStyle.Render("🧪 Bitcoin Transaction Builder") + "\n\n" +
				titleStyle.Render(m.errorMsg),
		)
	}

	header := titleStyle.Render("🧪 Bitcoin Transaction Builder")
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("────────────────────────────")

	content := lipgloss.JoinVertical(lipgloss.Left,
		renderStatusBox(m),
		renderInputBox(m),
	)

	return lipgloss.NewStyle().Render(header + "\n" + divider + "\n\n" + content)
}

func renderStatusBox(m model) string {
	toMap := formatToMap(m.toList, m.amountList)

	status := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("🌐 Network:     %s", m.net),
		fmt.Sprintf("📤 From:        %s", m.from),
		toMap,
		fmt.Sprintf("💰 Fund Addr:   %s", m.fundAddress),
	)

	return lipgloss.NewStyle().
		Width(90).
		Border(lipgloss.RoundedBorder()).
		Background(lipgloss.Color("234")).
		Padding(1, 1).
		Render(status)
}

func renderInputBox(m model) string {
	instructions := buildInputInstruction(m.step)

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 1).
		Width(80).
		Render(m.inputBuffer + "▌")

	inputHint := lipgloss.NewStyle().Italic(true).Render("Type input and press Enter ⏎")

	var preview string
	if m.step == 2 && m.inputBuffer != "" {
		preview = inputStyle.Render(fmt.Sprintf("➡️ Preview: %s", m.inputBuffer))
	}

	var errorView string
	if m.errorMsg != "" {
		errorView = errorStyle.Render("⚠️ " + m.errorMsg)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		instructions,
		"",
		inputBox,
		preview,
		inputHint,
		"",
		errorView,
	)

	return lipgloss.NewStyle().
		Width(90).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 1).
		Render(content)
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
		return highlightStyle.Render("Enter the Network:") + "\n(btc, btc-testnet3, btc-testnet4)"
	case 1:
		return highlightStyle.Render("Enter the 'From' address:")
	case 2:
		return highlightStyle.Render("Enter 'To' address and amount:") + "\n(e.g., bc1... 10000) or 'done'"
	case 3:
		return highlightStyle.Render("Enter the Fund address (optional):")
	default:
		return ""
	}
}
