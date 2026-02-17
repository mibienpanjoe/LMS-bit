package loan

import (
	"errors"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type Policy struct {
	LoanDays          int
	MaxLoansPerMember int
	MaxRenewals       int
}

func (p Policy) Validate() error {
	if p.LoanDays <= 0 {
		return errors.New("loan days must be greater than zero")
	}

	if p.MaxLoansPerMember <= 0 {
		return errors.New("max loans per member must be greater than zero")
	}

	if p.MaxRenewals < 0 {
		return errors.New("max renewals cannot be negative")
	}

	return nil
}

func CanIssue(c copy.Copy, m member.Member, activeLoans int, p Policy) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if !c.IsAvailable() {
		return shared.ErrCopyNotAvailable
	}

	if !m.CanBorrow() {
		return shared.ErrMemberNotEligible
	}

	if activeLoans >= p.MaxLoansPerMember {
		return shared.ErrLoanLimitReached
	}

	return nil
}

func New(id, copyID, memberID string, issuedAt time.Time, p Policy) (Loan, error) {
	if err := p.Validate(); err != nil {
		return Loan{}, err
	}

	l := Loan{
		ID:           id,
		CopyID:       copyID,
		MemberID:     memberID,
		IssuedAt:     issuedAt,
		DueAt:        issuedAt.AddDate(0, 0, p.LoanDays),
		RenewalCount: 0,
		Status:       StatusActive,
	}

	if err := l.Validate(); err != nil {
		return Loan{}, err
	}

	return l, nil
}

func Renew(l Loan, now time.Time, p Policy) (Loan, error) {
	if err := p.Validate(); err != nil {
		return Loan{}, err
	}

	if l.Status != StatusActive || l.ReturnedAt != nil {
		return Loan{}, shared.ErrLoanAlreadyClosed
	}

	if l.IsOverdue(now) {
		return Loan{}, shared.ErrLoanAlreadyOverdue
	}

	if l.RenewalCount >= p.MaxRenewals {
		return Loan{}, shared.ErrRenewalLimit
	}

	l.RenewalCount++
	l.DueAt = l.DueAt.AddDate(0, 0, p.LoanDays)

	return l, nil
}

func Return(l Loan, returnedAt time.Time) (Loan, error) {
	if l.Status == StatusReturned || l.ReturnedAt != nil {
		return Loan{}, shared.ErrLoanAlreadyClosed
	}

	l.Status = StatusReturned
	l.ReturnedAt = &returnedAt

	if err := l.Validate(); err != nil {
		return Loan{}, err
	}

	return l, nil
}
