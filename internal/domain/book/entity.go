package book

import (
	"errors"
	"regexp"
	"strings"
)

var isbnPattern = regexp.MustCompile(`^(?:\d{10}|\d{13})$`)

type Status string

const (
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
)

type Book struct {
	ID        string
	Title     string
	Authors   []string
	ISBN      string
	Category  string
	Publisher string
	Year      int
	Status    Status
}

func (b Book) Validate() error {
	if strings.TrimSpace(b.ID) == "" {
		return errors.New("book id is required")
	}

	if strings.TrimSpace(b.Title) == "" {
		return errors.New("book title is required")
	}

	if len(b.Authors) == 0 {
		return errors.New("at least one author is required")
	}

	if b.ISBN != "" && !isbnPattern.MatchString(strings.TrimSpace(b.ISBN)) {
		return errors.New("isbn must be 10 or 13 digits")
	}

	if b.Status == "" {
		return errors.New("book status is required")
	}

	return nil
}

func (b Book) CanCirculate() bool {
	return b.Status == StatusActive
}
