package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
)

func Info() string {

	strOut := strings.Builder{}

	tag := lipgloss.NewStyle().Align(lipgloss.Left).Render(time.Now().Format(time.Kitchen))
	info := lipgloss.NewStyle().Align(lipgloss.Right).MarginLeft(1).Foreground(green).Render("🍇 watching for changes ")
	ui := lipgloss.JoinVertical(lipgloss.Center, tag)
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	infoRender := infoBox.Render(ui)
	log := lipgloss.JoinHorizontal(lipgloss.Center, infoRender, info)
	strOut.WriteString(log)

	return strOut.String()
}
