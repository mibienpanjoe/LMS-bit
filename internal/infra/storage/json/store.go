package jsonstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
)

const schemaVersion = 1

var (
	ErrCorruptData      = errors.New("storage data is corrupt")
	ErrUnsupportedStore = errors.New("unsupported storage schema version")
)

type Store struct {
	mu   sync.RWMutex
	path string
	data snapshot
}

type snapshot struct {
	Version int                      `json:"version"`
	Books   map[string]book.Book     `json:"books"`
	Copies  map[string]copy.Copy     `json:"copies"`
	Members map[string]member.Member `json:"members"`
	Loans   map[string]loan.Loan     `json:"loans"`
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	s := &Store{path: path, data: newSnapshot()}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := s.writeSnapshot(s.data); err != nil {
			return nil, err
		}
		return s, nil
	} else if err != nil {
		return nil, fmt.Errorf("stat storage file: %w", err)
	}

	loaded, err := s.readSnapshot()
	if err != nil {
		return nil, err
	}

	s.data = loaded
	return s, nil
}

func newSnapshot() snapshot {
	return snapshot{
		Version: schemaVersion,
		Books:   map[string]book.Book{},
		Copies:  map[string]copy.Copy{},
		Members: map[string]member.Member{},
		Loans:   map[string]loan.Loan{},
	}
}

func (s *Store) readSnapshot() (snapshot, error) {
	content, err := os.ReadFile(s.path)
	if err != nil {
		return snapshot{}, fmt.Errorf("read storage file: %w", err)
	}

	if len(content) == 0 {
		return newSnapshot(), nil
	}

	var snap snapshot
	if err := json.Unmarshal(content, &snap); err != nil {
		return snapshot{}, fmt.Errorf("%w: decode json: %v", ErrCorruptData, err)
	}

	if snap.Version != schemaVersion {
		return snapshot{}, fmt.Errorf("%w: got %d expected %d", ErrUnsupportedStore, snap.Version, schemaVersion)
	}

	normalizeSnapshot(&snap)
	if err := validateSnapshot(snap); err != nil {
		return snapshot{}, err
	}

	return snap, nil
}

func (s *Store) writeSnapshot(snap snapshot) error {
	raw, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("encode storage json: %w", err)
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, raw, 0o644); err != nil {
		return fmt.Errorf("write temp storage file: %w", err)
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		return fmt.Errorf("atomic replace storage file: %w", err)
	}

	return nil
}

func normalizeSnapshot(s *snapshot) {
	if s.Books == nil {
		s.Books = map[string]book.Book{}
	}
	if s.Copies == nil {
		s.Copies = map[string]copy.Copy{}
	}
	if s.Members == nil {
		s.Members = map[string]member.Member{}
	}
	if s.Loans == nil {
		s.Loans = map[string]loan.Loan{}
	}
}

func validateSnapshot(s snapshot) error {
	for _, b := range s.Books {
		if err := b.Validate(); err != nil {
			return fmt.Errorf("%w: invalid book %q: %v", ErrCorruptData, b.ID, err)
		}
	}

	for _, c := range s.Copies {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("%w: invalid copy %q: %v", ErrCorruptData, c.ID, err)
		}
	}

	for _, m := range s.Members {
		if err := m.Validate(); err != nil {
			return fmt.Errorf("%w: invalid member %q: %v", ErrCorruptData, m.ID, err)
		}
	}

	for _, l := range s.Loans {
		if err := l.Validate(); err != nil {
			return fmt.Errorf("%w: invalid loan %q: %v", ErrCorruptData, l.ID, err)
		}
	}

	return nil
}
