package jsonstore

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type LoanRepository struct {
	store *Store
}

func NewLoanRepository(store *Store) *LoanRepository {
	return &LoanRepository{store: store}
}

func (r *LoanRepository) Save(_ context.Context, l loan.Loan) error {
	if err := l.Validate(); err != nil {
		return err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	r.store.data.Loans[l.ID] = l
	return r.store.writeSnapshot(r.store.data)
}

func (r *LoanRepository) GetByID(_ context.Context, id string) (loan.Loan, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	l, ok := r.store.data.Loans[id]
	if !ok {
		return loan.Loan{}, shared.ErrNotFound
	}

	return l, nil
}

func (r *LoanRepository) CountActiveByMemberID(_ context.Context, memberID string) (int, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	count := 0
	for _, l := range r.store.data.Loans {
		if l.MemberID == memberID && l.Status == loan.StatusActive {
			count++
		}
	}

	return count, nil
}

func (r *LoanRepository) List(_ context.Context) ([]loan.Loan, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]loan.Loan, 0, len(r.store.data.Loans))
	for _, l := range r.store.data.Loans {
		out = append(out, l)
	}

	return out, nil
}
