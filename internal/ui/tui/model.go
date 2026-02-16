package tui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
)

type Model struct {
	config   config.Config
	logger   *slog.Logger
	keys     keyMap
	help     help.Model
	styles   styles
	showHelp bool
	width    int
	height   int
}

func NewModel(cfg config.Config, logger *slog.Logger) Model {
	h := help.New()
	h.ShowAll = false

	return Model{
		config: cfg,
		logger: logger,
		keys:   newKeyMap(),
		help:   h,
		styles: newStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	m.logger.Info("starting tui", "app", m.config.AppName)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch {
		case keyMatches(msg.String(), m.keys.Quit.Keys()...):
			m.logger.Info("quitting")
			return m, tea.Quit
		case keyMatches(msg.String(), m.keys.Help.Keys()...):
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	header := m.styles.Header.Render(fmt.Sprintf(" %s - Phase 0 Bootstrap ", m.config.AppName))
	body := m.styles.Body.Render("Bubble Tea app is running. Next step: implement domain and use cases.")
	footer := m.styles.Footer.Render(m.help.View(m.keys))

	return strings.Join([]string{header, body, footer}, "\n")
}

func keyMatches(current string, expected ...string) bool {
	for _, e := range expected {
		if current == e {
			return true
		}
	}

	return false
}
