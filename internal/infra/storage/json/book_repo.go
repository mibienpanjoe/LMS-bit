package jsonstore

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type BookRepository struct {
	store *Store
}

func NewBookRepository(store *Store) *BookRepository {
	return &BookRepository{store: store}
}

func (r *BookRepository) Save(_ context.Context, b book.Book) error {
	if err := b.Validate(); err != nil {
		return err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	r.store.data.Books[b.ID] = b
	return r.store.writeSnapshot(r.store.data)
}

func (r *BookRepository) GetByID(_ context.Context, id string) (book.Book, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	b, ok := r.store.data.Books[id]
	if !ok {
		return book.Book{}, shared.ErrNotFound
	}

	return b, nil
}
