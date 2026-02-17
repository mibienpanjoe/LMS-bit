package loan

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusReturned Status = "returned"
)

type Loan struct {
	ID           string
	CopyID       string
	MemberID     string
	IssuedAt     time.Time
	DueAt        time.Time
	ReturnedAt   *time.Time
	RenewalCount int
	Status       Status
}

func (l Loan) Validate() error {
	if strings.TrimSpace(l.ID) == "" {
		return errors.New("loan id is required")
	}

	if strings.TrimSpace(l.CopyID) == "" {
		return errors.New("copy id is required")
	}

	if strings.TrimSpace(l.MemberID) == "" {
		return errors.New("member id is required")
	}

	if l.IssuedAt.IsZero() {
		return errors.New("issued date is required")
	}

	if l.DueAt.IsZero() || !l.DueAt.After(l.IssuedAt) {
		return errors.New("due date must be after issue date")
	}

	if l.Status == "" {
		return errors.New("loan status is required")
	}

	if l.Status == StatusReturned && l.ReturnedAt == nil {
		return errors.New("returned date is required when status is returned")
	}

	return nil
}

func (l Loan) IsOverdue(now time.Time) bool {
	if l.ReturnedAt != nil {
		return false
	}

	return now.After(l.DueAt)
}
