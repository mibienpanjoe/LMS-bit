package loan_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

func TestCanIssue(t *testing.T) {
	t.Parallel()

	p := loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1}

	tests := []struct {
		name   string
		copy   copy.Copy
		member member.Member
		active int
		want   error
	}{
		{
			name: "success",
			copy: copy.Copy{ID: "c-1", BookID: "b-1", Status: copy.StatusAvailable},
			member: member.Member{
				ID: "m-1", Name: "A", JoinedAt: time.Now(), Status: member.StatusActive,
			},
			active: 0,
			want:   nil,
		},
		{
			name: "copy unavailable",
			copy: copy.Copy{ID: "c-1", BookID: "b-1", Status: copy.StatusLoaned},
			member: member.Member{
				ID: "m-1", Name: "A", JoinedAt: time.Now(), Status: member.StatusActive,
			},
			active: 0,
			want:   shared.ErrCopyNotAvailable,
		},
		{
			name: "member blocked",
			copy: copy.Copy{ID: "c-1", BookID: "b-1", Status: copy.StatusAvailable},
			member: member.Member{
				ID: "m-1", Name: "A", JoinedAt: time.Now(), Status: member.StatusBlocked,
			},
			active: 0,
			want:   shared.ErrMemberNotEligible,
		},
		{
			name: "loan limit reached",
			copy: copy.Copy{ID: "c-1", BookID: "b-1", Status: copy.StatusAvailable},
			member: member.Member{
				ID: "m-1", Name: "A", JoinedAt: time.Now(), Status: member.StatusActive,
			},
			active: 3,
			want:   shared.ErrLoanLimitReached,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := loan.CanIssue(tc.copy, tc.member, tc.active, p)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v got %v", tc.want, err)
			}
		})
	}
}

func TestRenew(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	p := loan.Policy{LoanDays: 14, MaxLoansPerMember: 3, MaxRenewals: 1}

	tests := []struct {
		name string
		in   loan.Loan
		want error
	}{
		{
			name: "success",
			in: loan.Loan{
				ID: "l-1", CopyID: "c-1", MemberID: "m-1",
				IssuedAt: now.AddDate(0, 0, -5), DueAt: now.AddDate(0, 0, 5), Status: loan.StatusActive,
			},
			want: nil,
		},
		{
			name: "already returned",
			in: loan.Loan{
				ID: "l-1", CopyID: "c-1", MemberID: "m-1",
				IssuedAt: now.AddDate(0, 0, -10), DueAt: now.AddDate(0, 0, -2), Status: loan.StatusReturned,
			},
			want: shared.ErrLoanAlreadyClosed,
		},
		{
			name: "overdue",
			in: loan.Loan{
				ID: "l-1", CopyID: "c-1", MemberID: "m-1",
				IssuedAt: now.AddDate(0, 0, -20), DueAt: now.AddDate(0, 0, -1), Status: loan.StatusActive,
			},
			want: shared.ErrLoanAlreadyOverdue,
		},
		{
			name: "renewal limit",
			in: loan.Loan{
				ID: "l-1", CopyID: "c-1", MemberID: "m-1",
				IssuedAt: now.AddDate(0, 0, -5), DueAt: now.AddDate(0, 0, 5), Status: loan.StatusActive, RenewalCount: 1,
			},
			want: shared.ErrRenewalLimit,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := loan.Renew(tc.in, now, p)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v got %v", tc.want, err)
			}
			if tc.want == nil {
				if out.RenewalCount != tc.in.RenewalCount+1 {
					t.Fatalf("expected renewal count increment")
				}
			}
		})
	}
}

func TestReturn(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	in := loan.Loan{
		ID:       "l-1",
		CopyID:   "c-1",
		MemberID: "m-1",
		IssuedAt: now.AddDate(0, 0, -5),
		DueAt:    now.AddDate(0, 0, 5),
		Status:   loan.StatusActive,
	}

	out, err := loan.Return(in, now)
	if err != nil {
		t.Fatalf("expected nil error got %v", err)
	}

	if out.Status != loan.StatusReturned {
		t.Fatalf("expected status returned got %s", out.Status)
	}

	if out.ReturnedAt == nil || !out.ReturnedAt.Equal(now) {
		t.Fatalf("expected returned date to be set")
	}
}
