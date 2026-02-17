package integration_test

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	jsonstore "github.com/mibienpanjoe/LMS-bit/internal/infra/storage/json"
)

func TestIssueRenewReturnFlow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clock := &testClock{now: time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)}
	ids := &seqID{}

	services := newServices(t, filepath.Join(t.TempDir(), "flow.json"), clock, ids, loan.Policy{
		LoanDays:          14,
		MaxLoansPerMember: 3,
		MaxRenewals:       1,
	})

	b, err := services.books.Create(ctx, dto.CreateBookInput{Title: "Refactoring", Authors: []string{"M. Fowler"}, ISBN: "9780201485677"})
	if err != nil {
		t.Fatalf("create book: %v", err)
	}

	c, err := services.copies.Create(ctx, dto.CreateCopyInput{BookID: b.ID, Barcode: "CP-100"})
	if err != nil {
		t.Fatalf("create copy: %v", err)
	}

	m, err := services.members.Register(ctx, dto.RegisterMemberInput{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("register member: %v", err)
	}

	issued, err := services.loans.Issue(ctx, dto.IssueLoanInput{CopyID: c.ID, MemberID: m.ID})
	if err != nil {
		t.Fatalf("issue loan: %v", err)
	}

	renewed, err := services.loans.Renew(ctx, dto.RenewLoanInput{LoanID: issued.ID})
	if err != nil {
		t.Fatalf("renew loan: %v", err)
	}

	if !renewed.DueAt.After(issued.DueAt) {
		t.Fatalf("expected renewed due date after original due date")
	}

	returned, err := services.loans.Return(ctx, dto.ReturnLoanInput{LoanID: issued.ID})
	if err != nil {
		t.Fatalf("return loan: %v", err)
	}

	if returned.Status != loan.StatusReturned {
		t.Fatalf("expected returned status got %s", returned.Status)
	}

	allCopies, err := services.copies.List(ctx)
	if err != nil {
		t.Fatalf("list copies: %v", err)
	}

	status := ""
	for _, it := range allCopies {
		if it.ID == c.ID {
			status = string(it.Status)
			break
		}
	}
	if status != string(copy.StatusAvailable) {
		t.Fatalf("expected copy status available got %s", status)
	}
}

func TestOverdueAndPersistenceAfterReopen(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storePath := filepath.Join(t.TempDir(), "persist.json")
	clock := &testClock{now: time.Date(2026, 4, 1, 8, 0, 0, 0, time.UTC)}
	ids := &seqID{}

	services := newServices(t, storePath, clock, ids, loan.Policy{
		LoanDays:          1,
		MaxLoansPerMember: 3,
		MaxRenewals:       1,
	})

	b, _ := services.books.Create(ctx, dto.CreateBookInput{Title: "Go in Action", Authors: []string{"K. Kennedy"}, ISBN: "9781617291784"})
	c, _ := services.copies.Create(ctx, dto.CreateCopyInput{BookID: b.ID, Barcode: "CP-200"})
	m, _ := services.members.Register(ctx, dto.RegisterMemberInput{Name: "Bob", Email: "bob@example.com"})
	_, err := services.loans.Issue(ctx, dto.IssueLoanInput{CopyID: c.ID, MemberID: m.ID})
	if err != nil {
		t.Fatalf("issue loan: %v", err)
	}

	clock.now = clock.now.AddDate(0, 0, 2)
	overdue, err := services.loans.ListOverdue(ctx)
	if err != nil {
		t.Fatalf("list overdue: %v", err)
	}
	if len(overdue) != 1 {
		t.Fatalf("expected 1 overdue loan got %d", len(overdue))
	}

	reopened := newServices(t, storePath, clock, ids, loan.Policy{
		LoanDays:          1,
		MaxLoansPerMember: 3,
		MaxRenewals:       1,
	})

	books, err := reopened.books.List(ctx)
	if err != nil {
		t.Fatalf("list books after reopen: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("expected 1 book after reopen got %d", len(books))
	}

	overdueAfterReopen, err := reopened.loans.ListOverdue(ctx)
	if err != nil {
		t.Fatalf("list overdue after reopen: %v", err)
	}
	if len(overdueAfterReopen) != 1 {
		t.Fatalf("expected 1 overdue loan after reopen got %d", len(overdueAfterReopen))
	}
}

type services struct {
	books   usecase.BookService
	copies  usecase.CopyService
	members usecase.MemberService
	loans   usecase.LoanService
}

func newServices(t *testing.T, path string, clock *testClock, ids *seqID, policy loan.Policy) services {
	t.Helper()

	store, err := jsonstore.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}

	bookRepo := jsonstore.NewBookRepository(store)
	copyRepo := jsonstore.NewCopyRepository(store)
	memberRepo := jsonstore.NewMemberRepository(store)
	loanRepo := jsonstore.NewLoanRepository(store)

	return services{
		books:   usecase.NewBookService(bookRepo, ids),
		copies:  usecase.NewCopyService(copyRepo, ids),
		members: usecase.NewMemberService(memberRepo, ids, clock),
		loans:   usecase.NewLoanService(loanRepo, copyRepo, memberRepo, ids, clock, policy),
	}
}

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	return c.now
}

type seqID struct {
	n int
}

func (g *seqID) NewID() string {
	g.n++
	return "id-" + strconv.Itoa(g.n)
}
