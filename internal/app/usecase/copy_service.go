package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

type CopyService struct {
	copies ports.CopyRepository
	idGen  ports.IDGenerator
}

func NewCopyService(copies ports.CopyRepository, idGen ports.IDGenerator) CopyService {
	return CopyService{copies: copies, idGen: idGen}
}

func (s CopyService) Create(ctx context.Context, input dto.CreateCopyInput) (copy.Copy, error) {
	id := input.ID
	if id == "" {
		id = s.idGen.NewID()
	}

	if _, err := s.copies.GetByID(ctx, id); err == nil {
		return copy.Copy{}, shared.ErrDuplicateID
	} else if !errors.Is(err, shared.ErrNotFound) {
		return copy.Copy{}, err
	}

	if strings.TrimSpace(input.Barcode) != "" {
		if _, err := s.copies.GetByBarcode(ctx, input.Barcode); err == nil {
			return copy.Copy{}, shared.ErrDuplicateBarcode
		} else if !errors.Is(err, shared.ErrNotFound) {
			return copy.Copy{}, err
		}
	}

	c := copy.Copy{
		ID:            id,
		BookID:        input.BookID,
		Barcode:       input.Barcode,
		Status:        copy.StatusAvailable,
		ConditionNote: input.ConditionNote,
	}

	if err := c.Validate(); err != nil {
		return copy.Copy{}, err
	}

	if err := s.copies.Save(ctx, c); err != nil {
		return copy.Copy{}, err
	}

	return c, nil
}

func (s CopyService) Update(ctx context.Context, input dto.UpdateCopyInput) (copy.Copy, error) {
	c, err := s.copies.GetByID(ctx, input.ID)
	if err != nil {
		return copy.Copy{}, err
	}

	if strings.TrimSpace(input.Barcode) != "" && input.Barcode != c.Barcode {
		existing, err := s.copies.GetByBarcode(ctx, input.Barcode)
		if err == nil && existing.ID != c.ID {
			return copy.Copy{}, shared.ErrDuplicateBarcode
		}
		if err != nil && !errors.Is(err, shared.ErrNotFound) {
			return copy.Copy{}, err
		}
	}

	c.Barcode = input.Barcode
	c.ConditionNote = input.ConditionNote
	c.Status = copy.Status(strings.ToLower(strings.TrimSpace(input.Status)))

	if err := c.Validate(); err != nil {
		return copy.Copy{}, err
	}

	if err := s.copies.Save(ctx, c); err != nil {
		return copy.Copy{}, err
	}

	return c, nil
}

func (s CopyService) GetByID(ctx context.Context, id string) (copy.Copy, error) {
	return s.copies.GetByID(ctx, id)
}

func (s CopyService) List(ctx context.Context) ([]copy.Copy, error) {
	return s.copies.List(ctx)
}
