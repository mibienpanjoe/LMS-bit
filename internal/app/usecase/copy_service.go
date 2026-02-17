package usecase

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
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

func (s CopyService) List(ctx context.Context) ([]copy.Copy, error) {
	return s.copies.List(ctx)
}
