package tui

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
)

const (
	minWidth  = 80
	minHeight = 20
)

type Model struct {
	config config.Config
	logger *slog.Logger
	keys   keyMap
	styles styles
	help   help.Model
	route  route
	width  int
	height int

	table table.Model

	searchInput textinput.Model
	searching   bool
	searchQuery string

	showHelp bool
	status   statusMessage

	pendingConfirm bool
}

func NewModel(cfg config.Config, logger *slog.Logger) Model {
	h := help.New()
	h.ShowAll = false

	t := table.New(
		table.WithColumns(columnsForRoute(routeDashboard)),
		table.WithRows(rowsForRoute(routeDashboard, "")),
		table.WithFocused(true),
	)
	t.SetHeight(10)
	t.SetStyles(tableStyles())

	search := textinput.New()
	search.Placeholder = "type to filter"
	search.CharLimit = 64
	search.Prompt = ""

	return Model{
		config:         cfg,
		logger:         logger,
		keys:           newKeyMap(),
		styles:         newStyles(),
		help:           h,
		route:          routeDashboard,
		table:          t,
		searchInput:    search,
		searchQuery:    "",
		status:         statusMessage{text: "Ready", kind: statusInfo},
		pendingConfirm: false,
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
		m.resizeTable()
		return m, nil
	case clearStatusMsg:
		if m.status.kind != statusInfo {
			m.status = statusMessage{text: "Ready", kind: statusInfo}
		}
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			m.logger.Info("quitting")
			return m, tea.Quit
		}

		if key.Matches(msg, m.keys.ToggleHelp) {
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp
			return m, nil
		}

		if m.pendingConfirm {
			return m.updateConfirm(msg)
		}

		if m.searching {
			if key.Matches(msg, m.keys.Cancel) {
				m.searching = false
				m.searchInput.Blur()
				return m, nil
			}

			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.searchQuery = m.searchInput.Value()
			m.refreshRouteData()
			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keys.NextRoute):
			m.route = nextRoute(m.route)
			m.refreshRouteData()
			return m, m.setStatus(fmt.Sprintf("Switched to %s", m.route), statusInfo)
		case key.Matches(msg, m.keys.PrevRoute):
			m.route = prevRoute(m.route)
			m.refreshRouteData()
			return m, m.setStatus(fmt.Sprintf("Switched to %s", m.route), statusInfo)
		case key.Matches(msg, m.keys.Dashboard):
			m.route = routeDashboard
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Books):
			m.route = routeBooks
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Members):
			m.route = routeMembers
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Loans):
			m.route = routeLoans
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Reports):
			m.route = routeReports
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Settings):
			m.route = routeSettings
			m.refreshRouteData()
			return m, nil
		case key.Matches(msg, m.keys.Search):
			m.searching = true
			m.searchInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Danger):
			m.pendingConfirm = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Accept) {
		m.pendingConfirm = false
		return m, m.setStatus("Confirmed (placeholder action)", statusSuccess)
	}

	if key.Matches(msg, m.keys.Reject) || key.Matches(msg, m.keys.Cancel) {
		m.pendingConfirm = false
		return m, m.setStatus("Confirmation cancelled", statusInfo)
	}

	return m, nil
}

func (m *Model) setStatus(text string, kind statusKind) tea.Cmd {
	m.status = statusMessage{text: text, kind: kind}
	return clearStatusCmd(3 * time.Second)
}

func (m *Model) refreshRouteData() {
	m.table.SetColumns(columnsForRoute(m.route))
	m.table.SetRows(rowsForRoute(m.route, m.searchQuery))
	m.resizeTable()
}

func (m *Model) resizeTable() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	tableWidth := m.width - 2
	if tableWidth < 20 {
		tableWidth = 20
	}
	m.table.SetWidth(tableWidth)

	tableHeight := m.height - 10
	if tableHeight < 5 {
		tableHeight = 5
	}
	m.table.SetHeight(tableHeight)
}

func (m Model) View() string {
	if m.width > 0 && m.height > 0 && (m.width < minWidth || m.height < minHeight) {
		return m.styles.TooSmallScreen.Render(
			fmt.Sprintf("Terminal too small (%dx%d). Resize to at least %dx%d.", m.width, m.height, minWidth, minHeight),
		)
	}

	header := m.renderHeader()
	body := m.styles.Body.Render(m.table.View())
	footer := m.renderFooter()

	parts := []string{header, body}
	if m.pendingConfirm {
		parts = append(parts, m.renderConfirm())
	}
	parts = append(parts, footer)

	return strings.Join(parts, "\n")
}

func (m Model) renderHeader() string {
	title := m.styles.Header.Render(fmt.Sprintf(" %s ", m.config.AppName))
	subtitle := m.styles.HeaderMuted.Render("Phase 3 - TUI shell, routes, keymaps")
	nav := m.renderRouteTabs()
	searchLine := m.renderSearchLine()

	return strings.Join([]string{title + " " + subtitle, nav, searchLine}, "\n")
}

func (m Model) renderRouteTabs() string {
	items := make([]string, 0, len(allRoutes))
	for _, r := range allRoutes {
		label := string(r)
		if r == m.route {
			items = append(items, m.styles.ActiveTab.Render(label))
			continue
		}

		items = append(items, m.styles.Tab.Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, items...)
}

func (m Model) renderSearchLine() string {
	state := "inactive"
	if m.searching {
		state = "active"
	}

	prefix := m.styles.SearchLabel.Render("Search [/] (" + state + "):")
	return prefix + " " + m.searchInput.View()
}

func (m Model) renderFooter() string {
	statusStyle := m.styles.StatusInfo
	if m.status.kind == statusSuccess {
		statusStyle = m.styles.StatusSuccess
	}

	status := statusStyle.Render(m.status.text)
	helpView := m.help.View(m.keys)
	footer := lipgloss.JoinVertical(lipgloss.Left, status, m.styles.Footer.Render(helpView))
	return footer
}

func (m Model) renderConfirm() string {
	msg := strings.Join([]string{
		m.styles.ConfirmTitle.Render("Confirm action"),
		"This is a placeholder destructive action for phase scaffolding.",
		"Press y to confirm or n/esc to cancel.",
	}, "\n")

	box := m.styles.ConfirmBox.Render(msg)
	if m.width <= 0 {
		return box
	}

	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(box)
}

func tableStyles() table.Styles {
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("23"))
	ts.Selected = ts.Selected.
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("24")).
		Bold(true)
	return ts
}

func columnsForRoute(r route) []table.Column {
	switch r {
	case routeBooks:
		return []table.Column{{Title: "ID", Width: 10}, {Title: "Title", Width: 26}, {Title: "Author", Width: 18}, {Title: "ISBN", Width: 14}, {Title: "Status", Width: 10}}
	case routeMembers:
		return []table.Column{{Title: "ID", Width: 10}, {Title: "Name", Width: 20}, {Title: "Email", Width: 24}, {Title: "Phone", Width: 16}, {Title: "Status", Width: 10}}
	case routeLoans:
		return []table.Column{{Title: "LoanID", Width: 10}, {Title: "CopyID", Width: 10}, {Title: "MemberID", Width: 10}, {Title: "Due", Width: 14}, {Title: "State", Width: 10}}
	case routeReports:
		return []table.Column{{Title: "Metric", Width: 28}, {Title: "Value", Width: 12}, {Title: "Notes", Width: 30}}
	case routeSettings:
		return []table.Column{{Title: "Key", Width: 30}, {Title: "Value", Width: 30}, {Title: "Source", Width: 20}}
	default:
		return []table.Column{{Title: "Area", Width: 22}, {Title: "Status", Width: 12}, {Title: "Details", Width: 42}}
	}
}

func rowsForRoute(r route, query string) []table.Row {
	var rows []table.Row

	switch r {
	case routeBooks:
		rows = []table.Row{
			{"BK-100", "Clean Code", "R. Martin", "9780132350884", "active"},
			{"BK-101", "Design Patterns", "GoF", "9780201633610", "active"},
			{"BK-102", "Pragmatic Programmer", "Hunt/Thomas", "9780201616224", "archived"},
		}
	case routeMembers:
		rows = []table.Row{
			{"MB-200", "Asha Patel", "asha@example.com", "+1-555-1001", "active"},
			{"MB-201", "Noah Kim", "noah@example.com", "+1-555-1002", "blocked"},
			{"MB-202", "Maya Singh", "maya@example.com", "+1-555-1003", "active"},
		}
	case routeLoans:
		rows = []table.Row{
			{"LN-300", "CP-10", "MB-200", "2026-03-01", "active"},
			{"LN-301", "CP-12", "MB-202", "2026-02-14", "overdue"},
			{"LN-302", "CP-09", "MB-201", "2026-01-22", "returned"},
		}
	case routeReports:
		rows = []table.Row{
			{"Total catalog books", "124", "includes archived titles"},
			{"Active loans", "38", "updated from local snapshot"},
			{"Overdue loans", "5", "review in overdue report"},
		}
	case routeSettings:
		rows = []table.Row{
			{"loan.default_days", "14", "config"},
			{"loan.max_per_member", "3", "config"},
			{"ui.theme", "classic", "default"},
		}
	default:
		rows = []table.Row{
			{"Books", "ready", "catalog and inventory workflows"},
			{"Members", "ready", "member registration and lifecycle"},
			{"Loans", "ready", "issue, renew, return and overdue"},
		}
	}

	if query == "" {
		return rows
	}

	filtered := make([]table.Row, 0, len(rows))
	q := strings.ToLower(query)
	for _, row := range rows {
		joined := strings.ToLower(strings.Join(row, " "))
		if strings.Contains(joined, q) {
			filtered = append(filtered, row)
		}
	}

	if len(filtered) == 0 {
		return []table.Row{{"-", "No matches", "Try another search term"}}
	}

	return filtered
}
