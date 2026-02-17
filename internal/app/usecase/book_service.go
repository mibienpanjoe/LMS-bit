package usecase

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
)

type BookService struct {
	books ports.BookRepository
	idGen ports.IDGenerator
}

func NewBookService(books ports.BookRepository, idGen ports.IDGenerator) BookService {
	return BookService{books: books, idGen: idGen}
}

func (s BookService) Create(ctx context.Context, input dto.CreateBookInput) (book.Book, error) {
	id := input.ID
	if id == "" {
		id = s.idGen.NewID()
	}

	b := book.Book{
		ID:        id,
		Title:     input.Title,
		Authors:   input.Authors,
		ISBN:      input.ISBN,
		Category:  input.Category,
		Publisher: input.Publisher,
		Year:      input.Year,
		Status:    book.StatusActive,
	}

	if err := b.Validate(); err != nil {
		return book.Book{}, err
	}

	if err := s.books.Save(ctx, b); err != nil {
		return book.Book{}, err
	}

	return b, nil
}

func (s BookService) Archive(ctx context.Context, id string) (book.Book, error) {
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return book.Book{}, err
	}

	b.Status = book.StatusArchived
	if err := s.books.Save(ctx, b); err != nil {
		return book.Book{}, err
	}

	return b, nil
}
