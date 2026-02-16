package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	Header lipgloss.Style
	Body   lipgloss.Style
	Footer lipgloss.Style
}

func newStyles() styles {
	return styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("24")).
			Padding(0, 1),
		Body: lipgloss.NewStyle().
			Padding(1, 2),
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1),
	}
}
