package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/config"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/infra/id"
	jsonstore "github.com/mibienpanjoe/LMS-bit/internal/infra/storage/json"
	timeutil "github.com/mibienpanjoe/LMS-bit/internal/infra/time"
	"github.com/mibienpanjoe/LMS-bit/internal/logging"
	"github.com/mibienpanjoe/LMS-bit/internal/ui/tui"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	logger := logging.New(cfg.LogLevel)
	store, err := jsonstore.Open(cfg.StoragePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "storage open error: %v\n", err)
		os.Exit(1)
	}

	bookRepo := jsonstore.NewBookRepository(store)
	copyRepo := jsonstore.NewCopyRepository(store)
	memberRepo := jsonstore.NewMemberRepository(store)
	loanRepo := jsonstore.NewLoanRepository(store)

	idGen := id.NewGenerator()
	clock := timeutil.NewClock()

	bookService := usecase.NewBookService(bookRepo, idGen)
	copyService := usecase.NewCopyService(copyRepo, idGen)
	memberService := usecase.NewMemberService(memberRepo, idGen, clock)
	loanService := usecase.NewLoanService(
		loanRepo,
		copyRepo,
		memberRepo,
		idGen,
		clock,
		loan.Policy{
			LoanDays:          cfg.LoanDays,
			MaxLoansPerMember: cfg.MaxLoansPerUser,
			MaxRenewals:       cfg.MaxLoanRenewals,
		},
	)

	services := tui.Services{
		Books:   bookService,
		Copies:  copyService,
		Members: memberService,
		Loans:   loanService,
	}

	seedInitialData(context.Background(), services)

	program := tea.NewProgram(tui.NewModel(cfg, logger, services), tea.WithAltScreen())

	done := make(chan error, 1)
	go func() {
		_, err := program.Run()
		done <- err
	}()

	select {
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "application error: %v\n", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		program.Quit()
		err = <-done
		if err != nil && !errors.Is(err, tea.ErrProgramKilled) {
			fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
			os.Exit(1)
		}
	}
}

func seedInitialData(ctx context.Context, services tui.Services) {
	books, err := services.Books.List(ctx)
	if err != nil || len(books) > 0 {
		return
	}

	b, err := services.Books.Create(ctx, dto.CreateBookInput{
		Title:   "The Go Programming Language",
		Authors: []string{"Alan Donovan"},
		ISBN:    "9780134190440",
	})
	if err != nil {
		return
	}

	c, err := services.Copies.Create(ctx, dto.CreateCopyInput{BookID: b.ID, Barcode: "GO-CP-01"})
	if err != nil {
		return
	}

	_, _ = services.Copies.Create(ctx, dto.CreateCopyInput{BookID: b.ID, Barcode: "GO-CP-02"})

	m, err := services.Members.Register(ctx, dto.RegisterMemberInput{Name: "Demo User", Email: "demo@local"})
	if err != nil {
		return
	}

	_, _ = services.Loans.Issue(ctx, dto.IssueLoanInput{CopyID: c.ID, MemberID: m.ID})
}
