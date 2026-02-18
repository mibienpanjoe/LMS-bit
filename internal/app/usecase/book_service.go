package usecase

import (
	"context"
	"errors"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
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

	if _, err := s.books.GetByID(ctx, id); err == nil {
		return book.Book{}, shared.ErrDuplicateID
	} else if !errors.Is(err, shared.ErrNotFound) {
		return book.Book{}, err
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

func (s BookService) Update(ctx context.Context, input dto.UpdateBookInput) (book.Book, error) {
	b, err := s.books.GetByID(ctx, input.ID)
	if err != nil {
		return book.Book{}, err
	}

	b.Title = input.Title
	b.Authors = input.Authors
	b.ISBN = input.ISBN
	b.Category = input.Category
	b.Publisher = input.Publisher
	b.Year = input.Year

	if err := b.Validate(); err != nil {
		return book.Book{}, err
	}

	if err := s.books.Save(ctx, b); err != nil {
		return book.Book{}, err
	}

	return b, nil
}

func (s BookService) GetByID(ctx context.Context, id string) (book.Book, error) {
	return s.books.GetByID(ctx, id)
}

func (s BookService) SetStatus(ctx context.Context, id string, status book.Status) (book.Book, error) {
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return book.Book{}, err
	}

	b.Status = status
	if err := b.Validate(); err != nil {
		return book.Book{}, err
	}

	if err := s.books.Save(ctx, b); err != nil {
		return book.Book{}, err
	}

	return b, nil
}

func (s BookService) Archive(ctx context.Context, id string) (book.Book, error) {
	return s.SetStatus(ctx, id, book.StatusArchived)
}

func (s BookService) List(ctx context.Context) ([]book.Book, error) {
	return s.books.List(ctx)
}
