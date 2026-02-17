package usecase_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
)

func BenchmarkLoanServiceListOverdue10k(b *testing.B) {
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	loanData := make(map[string]loan.Loan, 10000)
	for i := 0; i < 10000; i++ {
		due := now.AddDate(0, 0, -1)
		if i%2 == 0 {
			due = now.AddDate(0, 0, 2)
		}

		id := "loan-" + strconv.Itoa(i)
		loanData[id] = loan.Loan{
			ID:       id,
			CopyID:   "copy-1",
			MemberID: "member-1",
			IssuedAt: now.AddDate(0, 0, -7),
			DueAt:    due,
			Status:   loan.StatusActive,
		}
	}

	svc := usecase.NewLoanService(
		&loanRepo{loans: loanData},
		&copyRepo{copies: map[string]copy.Copy{}},
		&memberRepo{members: map[string]member.Member{}},
		stubIDGen{id: "l-1"},
		stubClock{now: now},
		loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := svc.ListOverdue(context.Background())
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if len(result) == 0 {
			b.Fatalf("expected non-empty overdue result")
		}
	}
}
