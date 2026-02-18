package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/infra/id"
	jsonstore "github.com/mibienpanjoe/LMS-bit/internal/infra/storage/json"
	timeutil "github.com/mibienpanjoe/LMS-bit/internal/infra/time"
	"github.com/mibienpanjoe/LMS-bit/internal/logging"
)

func TestTabNavigationAcrossRoutesDoesNotPanic(t *testing.T) {
	t.Parallel()

	store, err := jsonstore.Open(filepath.Join(t.TempDir(), "storage.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}

	bookRepo := jsonstore.NewBookRepository(store)
	copyRepo := jsonstore.NewCopyRepository(store)
	memberRepo := jsonstore.NewMemberRepository(store)
	loanRepo := jsonstore.NewLoanRepository(store)

	idGen := id.NewGenerator()
	clock := timeutil.NewClock()

	services := Services{
		Books:   usecase.NewBookService(bookRepo, idGen),
		Copies:  usecase.NewCopyService(copyRepo, idGen),
		Members: usecase.NewMemberService(memberRepo, idGen, clock),
		Loans: usecase.NewLoanService(
			loanRepo,
			copyRepo,
			memberRepo,
			idGen,
			clock,
			loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1},
		),
	}

	cfg := config.Config{
		AppName:         "LMS-bit",
		LogLevel:        "error",
		StoragePath:     filepath.Join(t.TempDir(), "unused.json"),
		LoanDays:        14,
		MaxLoansPerUser: 3,
		MaxLoanRenewals: 1,
	}

	model := NewModel(cfg, logging.New("error"), services)

	for i := 0; i < 20; i++ {
		next, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
		model = next.(Model)
		_ = model.View()
	}
}
