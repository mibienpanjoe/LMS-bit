package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	Header         lipgloss.Style
	HeaderMuted    lipgloss.Style
	Tab            lipgloss.Style
	ActiveTab      lipgloss.Style
	Body           lipgloss.Style
	SearchLabel    lipgloss.Style
	Footer         lipgloss.Style
	StatusInfo     lipgloss.Style
	StatusSuccess  lipgloss.Style
	ConfirmBox     lipgloss.Style
	ConfirmTitle   lipgloss.Style
	TooSmallScreen lipgloss.Style
}

func newStyles() styles {
	borderColor := lipgloss.Color("238")

	return styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("24")).
			Padding(0, 1),
		HeaderMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Tab: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("249")),
		ActiveTab: lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("27")).
			Bold(true),
		Body: lipgloss.NewStyle().
			Padding(1, 1),
		SearchLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("110")),
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1),
		StatusInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("153")).
			Padding(0, 1),
		StatusSuccess: lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Padding(0, 1),
		ConfirmBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2),
		ConfirmTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("223")),
		TooSmallScreen: lipgloss.NewStyle().
			Foreground(lipgloss.Color("216")).
			Padding(1, 2),
	}
}
