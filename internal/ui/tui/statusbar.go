package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type statusKind string

const (
	statusInfo    statusKind = "info"
	statusSuccess statusKind = "success"
)

type statusMessage struct {
	text string
	kind statusKind
}

type clearStatusMsg struct{}

func clearStatusCmd(after time.Duration) tea.Cmd {
	return tea.Tick(after, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}
