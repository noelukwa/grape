package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// colors
	green  = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	wine   = lipgloss.AdaptiveColor{Light: "#C94C4C", Dark: "#F97F7F"}
	orange = lipgloss.AdaptiveColor{Light: "#FFA500", Dark: "#F28C28"}

	// text styles
	warn     = lipgloss.NewStyle().Foreground(orange).Bold(true).Render
	info     = lipgloss.NewStyle().Foreground(green).Render
	date     = lipgloss.NewStyle().Foreground(wine).SetString(time.Now().Format("15:04:05")).Padding(0, 1).String()
	infoText = func(s string) string {
		return info("[" + date + info("] ") + info(s))
	}

	delText = func(s string) string {

		return infoText(warn(">>>>> " + info(s)))

	}
)
