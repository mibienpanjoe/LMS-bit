package jsonstore_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	jsonstore "github.com/mibienpanjoe/LMS-bit/internal/infra/storage/json"
)

func TestStoreReadWriteReadConsistency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "storage.json")

	store, err := jsonstore.Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}

	bookRepo := jsonstore.NewBookRepository(store)
	copyRepo := jsonstore.NewCopyRepository(store)
	memberRepo := jsonstore.NewMemberRepository(store)
	loanRepo := jsonstore.NewLoanRepository(store)

	issuedAt := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)
	dueAt := issuedAt.AddDate(0, 0, 14)

	b := book.Book{
		ID:      "book-1",
		Title:   "Domain-Driven Design",
		Authors: []string{"Eric Evans"},
		ISBN:    "1234567890",
		Status:  book.StatusActive,
	}
	m := member.Member{
		ID:       "member-1",
		Name:     "Joe",
		JoinedAt: issuedAt,
		Status:   member.StatusActive,
	}
	c := copy.Copy{
		ID:     "copy-1",
		BookID: b.ID,
		Status: copy.StatusLoaned,
	}
	l := loan.Loan{
		ID:       "loan-1",
		CopyID:   c.ID,
		MemberID: m.ID,
		IssuedAt: issuedAt,
		DueAt:    dueAt,
		Status:   loan.StatusActive,
	}

	if err := bookRepo.Save(ctx, b); err != nil {
		t.Fatalf("save book: %v", err)
	}
	if err := memberRepo.Save(ctx, m); err != nil {
		t.Fatalf("save member: %v", err)
	}
	if err := copyRepo.Save(ctx, c); err != nil {
		t.Fatalf("save copy: %v", err)
	}
	if err := loanRepo.Save(ctx, l); err != nil {
		t.Fatalf("save loan: %v", err)
	}

	reopened, err := jsonstore.Open(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}

	bookRepo2 := jsonstore.NewBookRepository(reopened)
	copyRepo2 := jsonstore.NewCopyRepository(reopened)
	memberRepo2 := jsonstore.NewMemberRepository(reopened)
	loanRepo2 := jsonstore.NewLoanRepository(reopened)

	gotBook, err := bookRepo2.GetByID(ctx, b.ID)
	if err != nil || gotBook.Title != b.Title {
		t.Fatalf("book mismatch: %v %+v", err, gotBook)
	}

	gotMember, err := memberRepo2.GetByID(ctx, m.ID)
	if err != nil || gotMember.Name != m.Name {
		t.Fatalf("member mismatch: %v %+v", err, gotMember)
	}

	gotCopy, err := copyRepo2.GetByID(ctx, c.ID)
	if err != nil || gotCopy.Status != c.Status {
		t.Fatalf("copy mismatch: %v %+v", err, gotCopy)
	}

	gotLoan, err := loanRepo2.GetByID(ctx, l.ID)
	if err != nil || gotLoan.DueAt != l.DueAt {
		t.Fatalf("loan mismatch: %v %+v", err, gotLoan)
	}

	activeCount, err := loanRepo2.CountActiveByMemberID(ctx, m.ID)
	if err != nil {
		t.Fatalf("count active loans: %v", err)
	}
	if activeCount != 1 {
		t.Fatalf("expected one active loan got %d", activeCount)
	}
}

func TestOpenFailsOnCorruptJSON(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "storage.json")
	if err := os.WriteFile(path, []byte("{invalid-json"), 0o644); err != nil {
		t.Fatalf("seed corrupt file: %v", err)
	}

	_, err := jsonstore.Open(path)
	if !errors.Is(err, jsonstore.ErrCorruptData) {
		t.Fatalf("expected ErrCorruptData got %v", err)
	}
}
