package member

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBlocked  Status = "blocked"
)

type Member struct {
	ID       string
	Name     string
	Email    string
	Phone    string
	JoinedAt time.Time
	Status   Status
}

func (m Member) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("member id is required")
	}

	if strings.TrimSpace(m.Name) == "" {
		return errors.New("member name is required")
	}

	if m.JoinedAt.IsZero() {
		return errors.New("member join date is required")
	}

	if m.Status == "" {
		return errors.New("member status is required")
	}

	return nil
}

func (m Member) CanBorrow() bool {
	return m.Status == StatusActive
}
