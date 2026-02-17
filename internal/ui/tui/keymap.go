package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit       key.Binding
	ToggleHelp key.Binding
	NextRoute  key.Binding
	PrevRoute  key.Binding
	Dashboard  key.Binding
	Books      key.Binding
	Members    key.Binding
	Loans      key.Binding
	Reports    key.Binding
	Settings   key.Binding
	Search     key.Binding
	Cancel     key.Binding
	Danger     key.Binding
	Accept     key.Binding
	Reject     key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		NextRoute: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next view"),
		),
		PrevRoute: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev view"),
		),
		Dashboard: key.NewBinding(
			key.WithKeys("1", "d"),
			key.WithHelp("1", "dashboard"),
		),
		Books: key.NewBinding(
			key.WithKeys("2", "b"),
			key.WithHelp("2", "books"),
		),
		Members: key.NewBinding(
			key.WithKeys("3", "m"),
			key.WithHelp("3", "members"),
		),
		Loans: key.NewBinding(
			key.WithKeys("4", "l"),
			key.WithHelp("4", "loans"),
		),
		Reports: key.NewBinding(
			key.WithKeys("5", "r"),
			key.WithHelp("5", "reports"),
		),
		Settings: key.NewBinding(
			key.WithKeys("6", "s"),
			key.WithHelp("6", "settings"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Danger: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "confirm"),
		),
		Accept: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "confirm yes"),
		),
		Reject: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "confirm no"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextRoute, k.Search, k.ToggleHelp, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextRoute, k.PrevRoute, k.Search, k.Cancel},
		{k.Dashboard, k.Books, k.Members, k.Loans, k.Reports, k.Settings},
		{k.Danger, k.Accept, k.Reject, k.ToggleHelp, k.Quit},
	}
}
