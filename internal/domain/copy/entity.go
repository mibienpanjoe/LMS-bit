package copy

import (
	"errors"
	"strings"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusLoaned    Status = "loaned"
	StatusDamaged   Status = "damaged"
	StatusLost      Status = "lost"
)

type Copy struct {
	ID            string
	BookID        string
	Barcode       string
	Status        Status
	ConditionNote string
}

func (c Copy) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("copy id is required")
	}

	if strings.TrimSpace(c.BookID) == "" {
		return errors.New("book id is required")
	}

	if c.Status == "" {
		return errors.New("copy status is required")
	}

	return nil
}

func (c Copy) IsAvailable() bool {
	return c.Status == StatusAvailable
}
