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
	red    = lipgloss.AdaptiveColor{Light: "#D2042D", Dark: "#DC143C"}

	// text styles
	warn     = lipgloss.NewStyle().Foreground(orange).Bold(true).Render
	fail     = lipgloss.NewStyle().Foreground(red).Bold(true).Render
	info     = lipgloss.NewStyle().Foreground(green).Render
	date     = lipgloss.NewStyle().Foreground(wine).SetString(time.Now().Format("15:04:05")).Padding(0, 1).String()
	infoText = func(s string) string {
		return info("[" + date + info("] ") + info(s))
	}

	delText = func(s string) string {
		return infoText(warn(">>>>> " + info(s)))
	}

	stopText = func() string {
		return infoText(warn(StopNotice))
	}

	failText = func(s string) string {
		return infoText(fail(">>>>> " + fail(s)))
	}
)
