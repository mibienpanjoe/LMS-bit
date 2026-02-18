package jsonstore

import (
	"context"
	"strings"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type CopyRepository struct {
	store *Store
}

func NewCopyRepository(store *Store) *CopyRepository {
	return &CopyRepository{store: store}
}

func (r *CopyRepository) Save(_ context.Context, c copy.Copy) error {
	if err := c.Validate(); err != nil {
		return err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	r.store.data.Copies[c.ID] = c
	return r.store.writeSnapshot(r.store.data)
}

func (r *CopyRepository) GetByID(_ context.Context, id string) (copy.Copy, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	c, ok := r.store.data.Copies[id]
	if !ok {
		return copy.Copy{}, shared.ErrNotFound
	}

	return c, nil
}

func (r *CopyRepository) GetByBarcode(_ context.Context, barcode string) (copy.Copy, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	want := strings.TrimSpace(barcode)
	for _, c := range r.store.data.Copies {
		if strings.TrimSpace(c.Barcode) == want {
			return c, nil
		}
	}

	return copy.Copy{}, shared.ErrNotFound
}

func (r *CopyRepository) List(_ context.Context) ([]copy.Copy, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]copy.Copy, 0, len(r.store.data.Copies))
	for _, c := range r.store.data.Copies {
		out = append(out, c)
	}

	return out, nil
}
