package tui

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
)

const (
	minWidth                 = 88
	minHeight                = 22
	statusErrorPrefix        = "Error: "
	settingsSourceEnvDefault = "env/default"
)

type Services struct {
	Books   usecase.BookService
	Copies  usecase.CopyService
	Members usecase.MemberService
	Loans   usecase.LoanService
}

type loanFilter string

const (
	loanFilterAll      loanFilter = "all"
	loanFilterActive   loanFilter = "active"
	loanFilterOverdue  loanFilter = "overdue"
	loanFilterReturned loanFilter = "returned"
)

type formKind int

const (
	formBook formKind = iota
	formCopy
	formMember
	formIssueLoan
)

type formState struct {
	kind   formKind
	title  string
	fields []textinput.Model
	focus  int
}

type confirmAction int

const (
	confirmNone confirmAction = iota
	confirmArchiveBook
	confirmToggleMember
)

type Model struct {
	ctx      context.Context
	config   config.Config
	logger   *slog.Logger
	services Services

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

	loanFilter loanFilter

	activeForm *formState
	confirming bool
	confirmAct confirmAction
}

func NewModel(cfg config.Config, logger *slog.Logger, services Services) Model {
	h := help.New()
	h.ShowAll = false

	t := table.New(
		table.WithColumns([]table.Column{{Title: "Area", Width: 20}, {Title: "Status", Width: 12}, {Title: "Details", Width: 40}}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
	)
	t.SetHeight(10)
	t.SetStyles(tableStyles())

	search := textinput.New()
	search.Placeholder = "type to filter"
	search.CharLimit = 64
	search.Prompt = ""

	m := Model{
		ctx:         context.Background(),
		config:      cfg,
		logger:      logger,
		services:    services,
		keys:        newKeyMap(),
		styles:      newStyles(),
		help:        h,
		route:       routeDashboard,
		table:       t,
		searchInput: search,
		status:      statusMessage{text: "Ready", kind: statusInfo},
		loanFilter:  loanFilterAll,
	}
	m.refreshRouteData()
	return m
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
		return m.updateKeyMsg(msg)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) updateKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if handled, next, cmd := m.handleGlobalKeys(msg); handled {
		return next, cmd
	}

	if m.confirming {
		return m.updateConfirm(msg)
	}

	if m.activeForm != nil {
		return m.updateForm(msg)
	}

	if m.searching {
		return m.handleSearchKeys(msg)
	}

	if handled, next, cmd := m.handleRouteNavigationKeys(msg); handled {
		return next, cmd
	}

	if handled, next, cmd := m.handleRouteJumpKeys(msg); handled {
		return next, cmd
	}

	if handled, next, cmd := m.handleActionKeys(msg); handled {
		return next, cmd
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleGlobalKeys(msg tea.KeyMsg) (bool, Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Quit) {
		m.logger.Info("quitting")
		return true, m, tea.Quit
	}

	if key.Matches(msg, m.keys.ToggleHelp) {
		m.showHelp = !m.showHelp
		m.help.ShowAll = m.showHelp
		return true, m, nil
	}

	return false, m, nil
}

func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m Model) handleRouteNavigationKeys(msg tea.KeyMsg) (bool, Model, tea.Cmd) {
	if key.Matches(msg, m.keys.NextRoute) {
		m.route = nextRoute(m.route)
		m.refreshRouteData()
		return true, m, m.setStatus(fmt.Sprintf("Switched to %s", m.route), statusInfo)
	}

	if key.Matches(msg, m.keys.PrevRoute) {
		m.route = prevRoute(m.route)
		m.refreshRouteData()
		return true, m, m.setStatus(fmt.Sprintf("Switched to %s", m.route), statusInfo)
	}

	return false, m, nil
}

func (m Model) handleRouteJumpKeys(msg tea.KeyMsg) (bool, Model, tea.Cmd) {
	target, ok := m.routeTargetForKey(msg)
	if !ok {
		return false, m, nil
	}

	m.route = target
	m.refreshRouteData()
	return true, m, nil
}

func (m Model) routeTargetForKey(msg tea.KeyMsg) (route, bool) {
	switch {
	case key.Matches(msg, m.keys.Dashboard):
		return routeDashboard, true
	case key.Matches(msg, m.keys.Books):
		return routeBooks, true
	case key.Matches(msg, m.keys.Members):
		return routeMembers, true
	case key.Matches(msg, m.keys.Loans):
		return routeLoans, true
	case key.Matches(msg, m.keys.Reports):
		return routeReports, true
	case key.Matches(msg, m.keys.Settings):
		return routeSettings, true
	default:
		return "", false
	}
}

func (m Model) handleActionKeys(msg tea.KeyMsg) (bool, Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Search) {
		m.searching = true
		m.searchInput.Focus()
		return true, m, nil
	}

	if key.Matches(msg, m.keys.Add) {
		next, cmd := m.startAddFlow()
		return true, next.(Model), cmd
	}

	if key.Matches(msg, m.keys.CreateCopy) {
		if m.route == routeBooks {
			m.startCopyForm(m.selectedID())
		}
		return true, m, nil
	}

	if key.Matches(msg, m.keys.Archive) {
		next, cmd := m.startArchiveOrToggleConfirm()
		return true, next.(Model), cmd
	}

	if key.Matches(msg, m.keys.Filter) {
		if m.route == routeLoans {
			m.cycleLoanFilter()
			m.refreshRouteData()
			return true, m, m.setStatus("Loan filter: "+string(m.loanFilter), statusInfo)
		}
		return true, m, nil
	}

	if key.Matches(msg, m.keys.Issue) {
		if m.route == routeLoans {
			m.startIssueForm()
		}
		return true, m, nil
	}

	if key.Matches(msg, m.keys.Renew) {
		if m.route == routeLoans {
			next, cmd := m.renewSelectedLoan()
			return true, next.(Model), cmd
		}
		return true, m, nil
	}

	if key.Matches(msg, m.keys.Return) {
		if m.route == routeLoans {
			next, cmd := m.returnSelectedLoan()
			return true, next.(Model), cmd
		}
		return true, m, nil
	}

	return false, m, nil
}

func (m Model) View() string {
	if m.width > 0 && m.height > 0 && (m.width < minWidth || m.height < minHeight) {
		return m.styles.TooSmallScreen.Render(
			fmt.Sprintf("Terminal too small (%dx%d). Resize to at least %dx%d.", m.width, m.height, minWidth, minHeight),
		)
	}

	parts := []string{m.renderHeader(), m.styles.Body.Render(m.table.View())}
	if m.activeForm != nil {
		parts = append(parts, m.renderForm())
	}
	if m.confirming {
		parts = append(parts, m.renderConfirm())
	}
	parts = append(parts, m.renderFooter())
	return strings.Join(parts, "\n")
}

func (m Model) startAddFlow() (tea.Model, tea.Cmd) {
	switch m.route {
	case routeBooks:
		m.startBookForm()
	case routeMembers:
		m.startMemberForm()
	default:
		return m, nil
	}
	return m, nil
}

func (m *Model) startBookForm() {
	m.activeForm = newForm(formBook, "Add Book", []string{"Title", "Author", "ISBN"}, nil)
}

func (m *Model) startCopyForm(bookID string) {
	defaults := map[int]string{}
	if bookID != "" {
		defaults[0] = bookID
	}
	m.activeForm = newForm(formCopy, "Add Copy", []string{"Book ID", "Barcode", "Condition Note"}, defaults)
}

func (m *Model) startMemberForm() {
	m.activeForm = newForm(formMember, "Register Member", []string{"Name", "Email", "Phone"}, nil)
}

func (m *Model) startIssueForm() {
	m.activeForm = newForm(formIssueLoan, "Issue Loan", []string{"Copy ID", "Member ID"}, nil)
}

func newForm(kind formKind, title string, labels []string, defaults map[int]string) *formState {
	fields := make([]textinput.Model, 0, len(labels))
	for i, l := range labels {
		in := textinput.New()
		in.Placeholder = l
		in.Prompt = ""
		in.CharLimit = 128
		if val, ok := defaults[i]; ok {
			in.SetValue(val)
		}
		fields = append(fields, in)
	}
	fields[0].Focus()

	return &formState{kind: kind, title: title, fields: fields, focus: 0}
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Cancel) {
		m.activeForm = nil
		return m, m.setStatus("Form cancelled", statusInfo)
	}

	if msg.String() == "enter" {
		if m.activeForm.focus == len(m.activeForm.fields)-1 {
			return m.submitForm()
		}
		m.focusNextField()
		return m, nil
	}

	if msg.String() == "tab" {
		m.focusNextField()
		return m, nil
	}

	if msg.String() == "shift+tab" {
		m.focusPrevField()
		return m, nil
	}

	var cmd tea.Cmd
	idx := m.activeForm.focus
	m.activeForm.fields[idx], cmd = m.activeForm.fields[idx].Update(msg)
	return m, cmd
}

func (m *Model) focusNextField() {
	m.activeForm.fields[m.activeForm.focus].Blur()
	m.activeForm.focus = (m.activeForm.focus + 1) % len(m.activeForm.fields)
	m.activeForm.fields[m.activeForm.focus].Focus()
}

func (m *Model) focusPrevField() {
	m.activeForm.fields[m.activeForm.focus].Blur()
	m.activeForm.focus--
	if m.activeForm.focus < 0 {
		m.activeForm.focus = len(m.activeForm.fields) - 1
	}
	m.activeForm.fields[m.activeForm.focus].Focus()
}

func (m Model) submitForm() (tea.Model, tea.Cmd) {
	f := m.activeForm
	if f == nil {
		return m, nil
	}

	get := func(i int) string { return strings.TrimSpace(f.fields[i].Value()) }

	var err error
	switch f.kind {
	case formBook:
		_, err = m.services.Books.Create(m.ctx, dto.CreateBookInput{Title: get(0), Authors: []string{get(1)}, ISBN: get(2)})
	case formCopy:
		_, err = m.services.Copies.Create(m.ctx, dto.CreateCopyInput{BookID: get(0), Barcode: get(1), ConditionNote: get(2)})
	case formMember:
		_, err = m.services.Members.Register(m.ctx, dto.RegisterMemberInput{Name: get(0), Email: get(1), Phone: get(2)})
	case formIssueLoan:
		_, err = m.services.Loans.Issue(m.ctx, dto.IssueLoanInput{CopyID: get(0), MemberID: get(1)})
	}

	if err != nil {
		return m, m.setStatus(statusErrorPrefix+err.Error(), statusInfo)
	}

	m.activeForm = nil
	m.refreshRouteData()
	return m, m.setStatus("Saved successfully", statusSuccess)
}

func (m Model) startArchiveOrToggleConfirm() (tea.Model, tea.Cmd) {
	id := m.selectedID()
	if id == "" {
		return m, m.setStatus("Select a row first", statusInfo)
	}

	switch m.route {
	case routeBooks:
		m.confirming = true
		m.confirmAct = confirmArchiveBook
		return m, nil
	case routeMembers:
		m.confirming = true
		m.confirmAct = confirmToggleMember
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, m.keys.Reject) || key.Matches(msg, m.keys.Cancel) {
		m.confirming = false
		m.confirmAct = confirmNone
		return m, m.setStatus("Confirmation cancelled", statusInfo)
	}

	if !key.Matches(msg, m.keys.Accept) {
		return m, nil
	}

	id := m.selectedID()
	err := m.performConfirmAction(id)

	m.confirming = false
	m.confirmAct = confirmNone
	if err != nil {
		return m, m.setStatus(statusErrorPrefix+err.Error(), statusInfo)
	}

	m.refreshRouteData()
	return m, m.setStatus("Updated successfully", statusSuccess)
}

func (m Model) performConfirmAction(id string) error {
	switch m.confirmAct {
	case confirmArchiveBook:
		_, err := m.services.Books.Archive(m.ctx, id)
		return err
	case confirmToggleMember:
		return m.toggleMemberStatus(id)
	default:
		return nil
	}
}

func (m Model) toggleMemberStatus(id string) error {
	members, err := m.services.Members.List(m.ctx)
	if err != nil {
		return err
	}

	for _, it := range members {
		if it.ID != id {
			continue
		}

		next := member.StatusInactive
		if it.Status != member.StatusActive {
			next = member.StatusActive
		}
		_, err = m.services.Members.SetStatus(m.ctx, id, next)
		return err
	}

	return nil
}

func (m Model) renewSelectedLoan() (tea.Model, tea.Cmd) {
	id := m.selectedID()
	if id == "" {
		return m, m.setStatus("Select a loan first", statusInfo)
	}

	if _, err := m.services.Loans.Renew(m.ctx, dto.RenewLoanInput{LoanID: id}); err != nil {
		return m, m.setStatus(statusErrorPrefix+err.Error(), statusInfo)
	}

	m.refreshRouteData()
	return m, m.setStatus("Loan renewed", statusSuccess)
}

func (m Model) returnSelectedLoan() (tea.Model, tea.Cmd) {
	id := m.selectedID()
	if id == "" {
		return m, m.setStatus("Select a loan first", statusInfo)
	}

	if _, err := m.services.Loans.Return(m.ctx, dto.ReturnLoanInput{LoanID: id}); err != nil {
		return m, m.setStatus(statusErrorPrefix+err.Error(), statusInfo)
	}

	m.refreshRouteData()
	return m, m.setStatus("Loan returned", statusSuccess)
}

func (m *Model) cycleLoanFilter() {
	switch m.loanFilter {
	case loanFilterAll:
		m.loanFilter = loanFilterActive
	case loanFilterActive:
		m.loanFilter = loanFilterOverdue
	case loanFilterOverdue:
		m.loanFilter = loanFilterReturned
	default:
		m.loanFilter = loanFilterAll
	}
}

func (m *Model) refreshRouteData() {
	var (
		cols []table.Column
		rows []table.Row
	)

	switch m.route {
	case routeBooks:
		cols, rows = m.booksTable()
	case routeMembers:
		cols, rows = m.membersTable()
	case routeLoans:
		cols, rows = m.loansTable()
	case routeReports:
		cols, rows = m.reportsTable()
	case routeSettings:
		cols, rows = m.settingsTable()
	default:
		cols, rows = m.dashboardTable()
	}

	filtered := filterRows(rows, m.searchQuery)
	m.table.SetColumns(cols)
	m.table.SetRows(filtered)
	m.resizeTable()
}

func (m Model) dashboardTable() ([]table.Column, []table.Row) {
	books, _ := m.services.Books.List(m.ctx)
	copies, _ := m.services.Copies.List(m.ctx)
	members, _ := m.services.Members.List(m.ctx)
	loans, _ := m.services.Loans.List(m.ctx)
	overdue, _ := m.services.Loans.ListOverdue(m.ctx)

	activeLoans := 0
	for _, l := range loans {
		if l.Status == loan.StatusActive {
			activeLoans++
		}
	}

	rows := []table.Row{
		{"Books", "ok", fmt.Sprintf("%d titles", len(books))},
		{"Copies", "ok", fmt.Sprintf("%d copies", len(copies))},
		{"Members", "ok", fmt.Sprintf("%d registered", len(members))},
		{"Active Loans", "ok", fmt.Sprintf("%d open", activeLoans)},
		{"Overdue", "watch", fmt.Sprintf("%d overdue", len(overdue))},
	}

	return []table.Column{{Title: "Area", Width: 20}, {Title: "Status", Width: 12}, {Title: "Details", Width: 40}}, rows
}

func (m Model) booksTable() ([]table.Column, []table.Row) {
	books, _ := m.services.Books.List(m.ctx)
	copies, _ := m.services.Copies.List(m.ctx)
	copyCount := map[string]int{}
	for _, c := range copies {
		copyCount[c.BookID]++
	}

	sort.Slice(books, func(i, j int) bool { return strings.ToLower(books[i].Title) < strings.ToLower(books[j].Title) })

	rows := make([]table.Row, 0, len(books))
	for _, b := range books {
		author := ""
		if len(b.Authors) > 0 {
			author = b.Authors[0]
		}
		rows = append(rows, table.Row{b.ID, b.Title, author, b.ISBN, fmt.Sprintf("%d", copyCount[b.ID]), string(b.Status)})
	}

	if len(rows) == 0 {
		rows = []table.Row{{"-", "No books yet", "Press a to add", "", "", ""}}
	}

	return []table.Column{{Title: "ID", Width: 12}, {Title: "Title", Width: 22}, {Title: "Author", Width: 16}, {Title: "ISBN", Width: 14}, {Title: "Copies", Width: 8}, {Title: "Status", Width: 10}}, rows
}

func (m Model) membersTable() ([]table.Column, []table.Row) {
	members, _ := m.services.Members.List(m.ctx)
	sort.Slice(members, func(i, j int) bool { return strings.ToLower(members[i].Name) < strings.ToLower(members[j].Name) })

	rows := make([]table.Row, 0, len(members))
	for _, mm := range members {
		rows = append(rows, table.Row{mm.ID, mm.Name, mm.Email, mm.Phone, string(mm.Status)})
	}

	if len(rows) == 0 {
		rows = []table.Row{{"-", "No members yet", "Press a to register", "", ""}}
	}

	return []table.Column{{Title: "ID", Width: 12}, {Title: "Name", Width: 20}, {Title: "Email", Width: 24}, {Title: "Phone", Width: 16}, {Title: "Status", Width: 10}}, rows
}

func (m Model) loansTable() ([]table.Column, []table.Row) {
	loans, _ := m.services.Loans.List(m.ctx)
	now := time.Now().UTC()
	sort.Slice(loans, func(i, j int) bool { return loans[i].DueAt.Before(loans[j].DueAt) })

	rows := make([]table.Row, 0, len(loans))
	for _, l := range loans {
		state := string(l.Status)
		if l.IsOverdue(now) {
			state = "overdue"
		}

		if !loanMatchesFilter(l, state, m.loanFilter) {
			continue
		}

		rows = append(rows, table.Row{l.ID, l.CopyID, l.MemberID, l.DueAt.Format("2006-01-02"), state})
	}

	if len(rows) == 0 {
		rows = []table.Row{{"-", "No loans", "Press i to issue", "", string(m.loanFilter)}}
	}

	return []table.Column{{Title: "LoanID", Width: 12}, {Title: "CopyID", Width: 12}, {Title: "MemberID", Width: 12}, {Title: "Due", Width: 14}, {Title: "State", Width: 12}}, rows
}

func (m Model) reportsTable() ([]table.Column, []table.Row) {
	overdue, _ := m.services.Loans.ListOverdue(m.ctx)
	sort.Slice(overdue, func(i, j int) bool { return overdue[i].DueAt.Before(overdue[j].DueAt) })

	rows := make([]table.Row, 0, len(overdue))
	for _, l := range overdue {
		rows = append(rows, table.Row{l.ID, l.MemberID, l.CopyID, l.DueAt.Format("2006-01-02")})
	}

	if len(rows) == 0 {
		rows = []table.Row{{"-", "No overdue loans", "", ""}}
	}

	return []table.Column{{Title: "LoanID", Width: 12}, {Title: "MemberID", Width: 12}, {Title: "CopyID", Width: 12}, {Title: "Due Date", Width: 14}}, rows
}

func (m Model) settingsTable() ([]table.Column, []table.Row) {
	rows := []table.Row{
		{"storage.path", m.config.StoragePath, settingsSourceEnvDefault},
		{"loan.days", fmt.Sprintf("%d", m.config.LoanDays), settingsSourceEnvDefault},
		{"loan.max_per_member", fmt.Sprintf("%d", m.config.MaxLoansPerUser), settingsSourceEnvDefault},
		{"loan.max_renewals", fmt.Sprintf("%d", m.config.MaxLoanRenewals), settingsSourceEnvDefault},
	}
	return []table.Column{{Title: "Key", Width: 28}, {Title: "Value", Width: 34}, {Title: "Source", Width: 18}}, rows
}

func loanMatchesFilter(l loan.Loan, state string, filter loanFilter) bool {
	switch filter {
	case loanFilterActive:
		return l.Status == loan.StatusActive && state != "overdue"
	case loanFilterOverdue:
		return state == "overdue"
	case loanFilterReturned:
		return l.Status == loan.StatusReturned
	default:
		return true
	}
}

func filterRows(rows []table.Row, query string) []table.Row {
	if query == "" {
		return rows
	}

	q := strings.ToLower(query)
	filtered := make([]table.Row, 0, len(rows))
	for _, row := range rows {
		if strings.Contains(strings.ToLower(strings.Join(row, " ")), q) {
			filtered = append(filtered, row)
		}
	}

	if len(filtered) == 0 {
		return []table.Row{{"-", "No matches", "Try another search term"}}
	}

	return filtered
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

	tableHeight := m.height - 11
	if m.activeForm != nil || m.confirming {
		tableHeight -= 6
	}
	if tableHeight < 5 {
		tableHeight = 5
	}
	m.table.SetHeight(tableHeight)
}

func (m Model) renderHeader() string {
	title := m.styles.Header.Render(fmt.Sprintf(" %s ", m.config.AppName))
	subtitle := m.styles.HeaderMuted.Render("Phase 4 - Feature views and workflows")
	nav := m.renderRouteTabs()
	searchLine := m.renderSearchLine()

	return strings.Join([]string{title + " " + subtitle, nav, searchLine}, "\n")
}

func (m Model) renderRouteTabs() string {
	items := make([]string, 0, len(allRoutes))
	for _, r := range allRoutes {
		label := string(r)
		if r == routeLoans {
			label = label + " (" + string(m.loanFilter) + ")"
		}
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

	helpView := m.help.View(m.keys)
	status := statusStyle.Render(m.status.text)
	return lipgloss.JoinVertical(lipgloss.Left, status, m.styles.Footer.Render(helpView))
}

func (m Model) renderForm() string {
	if m.activeForm == nil {
		return ""
	}

	lines := []string{m.styles.ConfirmTitle.Render(m.activeForm.title)}
	for i, f := range m.activeForm.fields {
		prefix := "  "
		if i == m.activeForm.focus {
			prefix = "> "
		}
		lines = append(lines, prefix+f.View())
	}
	lines = append(lines, "tab/shift+tab to move, enter to continue/submit, esc to cancel")

	box := m.styles.ConfirmBox.Render(strings.Join(lines, "\n"))
	if m.width <= 0 {
		return box
	}
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(box)
}

func (m Model) renderConfirm() string {
	title := "Confirm"
	body := "Proceed with this action?"
	if m.confirmAct == confirmArchiveBook {
		title = "Archive Book"
		body = "This will mark the selected book as archived."
	}
	if m.confirmAct == confirmToggleMember {
		title = "Toggle Member Status"
		body = "This will toggle selected member active/inactive status."
	}

	msg := strings.Join([]string{
		m.styles.ConfirmTitle.Render(title),
		body,
		"Press y to confirm or n/esc to cancel.",
	}, "\n")

	box := m.styles.ConfirmBox.Render(msg)
	if m.width <= 0 {
		return box
	}
	return lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(box)
}

func (m *Model) setStatus(text string, kind statusKind) tea.Cmd {
	m.status = statusMessage{text: text, kind: kind}
	return clearStatusCmd(3 * time.Second)
}

func (m Model) selectedID() string {
	row := m.table.SelectedRow()
	if len(row) == 0 {
		return ""
	}
	if row[0] == "-" {
		return ""
	}
	return row[0]
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
