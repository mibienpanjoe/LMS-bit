package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

func TestLoanServiceIssue(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	copyRepo := &copyRepo{copies: map[string]copy.Copy{
		"c-1": {ID: "c-1", BookID: "b-1", Status: copy.StatusAvailable},
	}}
	memberRepo := &memberRepo{members: map[string]member.Member{
		"m-1": {ID: "m-1", Name: "Joe", JoinedAt: now, Status: member.StatusActive},
	}}
	loanRepo := &loanRepo{loans: map[string]loan.Loan{}}

	svc := usecase.NewLoanService(
		loanRepo,
		copyRepo,
		memberRepo,
		stubIDGen{id: "l-1"},
		stubClock{now: now},
		loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1},
	)

	issued, err := svc.Issue(context.Background(), dto.IssueLoanInput{CopyID: "c-1", MemberID: "m-1"})
	if err != nil {
		t.Fatalf("expected nil error got %v", err)
	}

	if issued.ID != "l-1" {
		t.Fatalf("expected loan id l-1 got %s", issued.ID)
	}

	storedCopy, err := copyRepo.GetByID(context.Background(), "c-1")
	if err != nil {
		t.Fatalf("expected copy in repo got %v", err)
	}

	if storedCopy.Status != copy.StatusLoaned {
		t.Fatalf("expected copy to be loaned got %s", storedCopy.Status)
	}
}

func TestLoanServiceIssueFailsWhenBlocked(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	copyRepo := &copyRepo{copies: map[string]copy.Copy{
		"c-1": {ID: "c-1", BookID: "b-1", Status: copy.StatusAvailable},
	}}
	memberRepo := &memberRepo{members: map[string]member.Member{
		"m-1": {ID: "m-1", Name: "Joe", JoinedAt: now, Status: member.StatusBlocked},
	}}
	loanRepo := &loanRepo{loans: map[string]loan.Loan{}}

	svc := usecase.NewLoanService(
		loanRepo,
		copyRepo,
		memberRepo,
		stubIDGen{id: "l-1"},
		stubClock{now: now},
		loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1},
	)

	_, err := svc.Issue(context.Background(), dto.IssueLoanInput{CopyID: "c-1", MemberID: "m-1"})
	if !errors.Is(err, shared.ErrMemberNotEligible) {
		t.Fatalf("expected %v got %v", shared.ErrMemberNotEligible, err)
	}
}

type stubClock struct {
	now time.Time
}

func (s stubClock) Now() time.Time {
	return s.now
}

type stubIDGen struct {
	id string
}

func (s stubIDGen) NewID() string {
	return s.id
}

type copyRepo struct {
	copies map[string]copy.Copy
}

func (r *copyRepo) Save(_ context.Context, c copy.Copy) error {
	r.copies[c.ID] = c
	return nil
}

func (r *copyRepo) GetByID(_ context.Context, id string) (copy.Copy, error) {
	c, ok := r.copies[id]
	if !ok {
		return copy.Copy{}, shared.ErrNotFound
	}
	return c, nil
}

func (r *copyRepo) List(_ context.Context) ([]copy.Copy, error) {
	out := make([]copy.Copy, 0, len(r.copies))
	for _, c := range r.copies {
		out = append(out, c)
	}
	return out, nil
}

type memberRepo struct {
	members map[string]member.Member
}

func (r *memberRepo) Save(_ context.Context, m member.Member) error {
	r.members[m.ID] = m
	return nil
}

func (r *memberRepo) GetByID(_ context.Context, id string) (member.Member, error) {
	m, ok := r.members[id]
	if !ok {
		return member.Member{}, shared.ErrNotFound
	}
	return m, nil
}

func (r *memberRepo) List(_ context.Context) ([]member.Member, error) {
	out := make([]member.Member, 0, len(r.members))
	for _, m := range r.members {
		out = append(out, m)
	}
	return out, nil
}

type loanRepo struct {
	loans map[string]loan.Loan
}

func (r *loanRepo) Save(_ context.Context, l loan.Loan) error {
	r.loans[l.ID] = l
	return nil
}

func (r *loanRepo) GetByID(_ context.Context, id string) (loan.Loan, error) {
	l, ok := r.loans[id]
	if !ok {
		return loan.Loan{}, shared.ErrNotFound
	}
	return l, nil
}

func (r *loanRepo) CountActiveByMemberID(_ context.Context, memberID string) (int, error) {
	count := 0
	for _, l := range r.loans {
		if l.MemberID == memberID && l.Status == loan.StatusActive {
			count++
		}
	}

	return count, nil
}

func (r *loanRepo) List(_ context.Context) ([]loan.Loan, error) {
	out := make([]loan.Loan, 0, len(r.loans))
	for _, l := range r.loans {
		out = append(out, l)
	}
	return out, nil
}

type bookRepo struct {
	books map[string]book.Book
}

func (r *bookRepo) Save(_ context.Context, b book.Book) error {
	r.books[b.ID] = b
	return nil
}

func (r *bookRepo) GetByID(_ context.Context, id string) (book.Book, error) {
	b, ok := r.books[id]
	if !ok {
		return book.Book{}, shared.ErrNotFound
	}
	return b, nil
}

func (r *bookRepo) List(_ context.Context) ([]book.Book, error) {
	out := make([]book.Book, 0, len(r.books))
	for _, b := range r.books {
		out = append(out, b)
	}
	return out, nil
}
