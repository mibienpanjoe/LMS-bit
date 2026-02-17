package usecase

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
)

type LoanService struct {
	loans   ports.LoanRepository
	copies  ports.CopyRepository
	members ports.MemberRepository
	idGen   ports.IDGenerator
	clock   ports.Clock
	policy  loan.Policy
}

func NewLoanService(
	loans ports.LoanRepository,
	copies ports.CopyRepository,
	members ports.MemberRepository,
	idGen ports.IDGenerator,
	clock ports.Clock,
	policy loan.Policy,
) LoanService {
	return LoanService{
		loans:   loans,
		copies:  copies,
		members: members,
		idGen:   idGen,
		clock:   clock,
		policy:  policy,
	}
}

func (s LoanService) Issue(ctx context.Context, input dto.IssueLoanInput) (loan.Loan, error) {
	c, err := s.copies.GetByID(ctx, input.CopyID)
	if err != nil {
		return loan.Loan{}, err
	}

	m, err := s.members.GetByID(ctx, input.MemberID)
	if err != nil {
		return loan.Loan{}, err
	}

	activeCount, err := s.loans.CountActiveByMemberID(ctx, input.MemberID)
	if err != nil {
		return loan.Loan{}, err
	}

	if err := loan.CanIssue(c, m, activeCount, s.policy); err != nil {
		return loan.Loan{}, err
	}

	created, err := loan.New(s.idGen.NewID(), input.CopyID, input.MemberID, s.clock.Now(), s.policy)
	if err != nil {
		return loan.Loan{}, err
	}

	if err := s.loans.Save(ctx, created); err != nil {
		return loan.Loan{}, err
	}

	c.Status = copy.StatusLoaned
	if err := s.copies.Save(ctx, c); err != nil {
		return loan.Loan{}, err
	}

	return created, nil
}

func (s LoanService) Renew(ctx context.Context, input dto.RenewLoanInput) (loan.Loan, error) {
	current, err := s.loans.GetByID(ctx, input.LoanID)
	if err != nil {
		return loan.Loan{}, err
	}

	renewed, err := loan.Renew(current, s.clock.Now(), s.policy)
	if err != nil {
		return loan.Loan{}, err
	}

	if err := s.loans.Save(ctx, renewed); err != nil {
		return loan.Loan{}, err
	}

	return renewed, nil
}

func (s LoanService) Return(ctx context.Context, input dto.ReturnLoanInput) (loan.Loan, error) {
	current, err := s.loans.GetByID(ctx, input.LoanID)
	if err != nil {
		return loan.Loan{}, err
	}

	returned, err := loan.Return(current, s.clock.Now())
	if err != nil {
		return loan.Loan{}, err
	}

	if err := s.loans.Save(ctx, returned); err != nil {
		return loan.Loan{}, err
	}

	c, err := s.copies.GetByID(ctx, returned.CopyID)
	if err != nil {
		return loan.Loan{}, err
	}

	c.Status = copy.StatusAvailable
	if err := s.copies.Save(ctx, c); err != nil {
		return loan.Loan{}, err
	}

	return returned, nil
}
