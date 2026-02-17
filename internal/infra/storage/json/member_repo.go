package jsonstore

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type MemberRepository struct {
	store *Store
}

func NewMemberRepository(store *Store) *MemberRepository {
	return &MemberRepository{store: store}
}

func (r *MemberRepository) Save(_ context.Context, m member.Member) error {
	if err := m.Validate(); err != nil {
		return err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	r.store.data.Members[m.ID] = m
	return r.store.writeSnapshot(r.store.data)
}

func (r *MemberRepository) GetByID(_ context.Context, id string) (member.Member, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	m, ok := r.store.data.Members[id]
	if !ok {
		return member.Member{}, shared.ErrNotFound
	}

	return m, nil
}

func (r *MemberRepository) List(_ context.Context) ([]member.Member, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]member.Member, 0, len(r.store.data.Members))
	for _, m := range r.store.data.Members {
		out = append(out, m)
	}

	return out, nil
}
