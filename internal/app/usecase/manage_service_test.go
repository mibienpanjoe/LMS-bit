package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/usecase"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/shared"
)

func TestBookServiceUpdate(t *testing.T) {
	t.Parallel()

	repo := &bookRepo{books: map[string]book.Book{
		"b-1": {ID: "b-1", Title: "Old", Authors: []string{"Author"}, Status: book.StatusActive},
	}}

	svc := usecase.NewBookService(repo, stubIDGen{id: "ignored"})

	updated, err := svc.Update(context.Background(), dto.UpdateBookInput{
		ID:        "b-1",
		Title:     "New Title",
		Authors:   []string{"New Author"},
		ISBN:      "1234567890",
		Category:  "Tech",
		Publisher: "Pub",
		Year:      2026,
	})
	if err != nil {
		t.Fatalf("update book: %v", err)
	}

	if updated.Title != "New Title" || updated.Category != "Tech" {
		t.Fatalf("unexpected updated book: %+v", updated)
	}
}

func TestBookServiceSetStatus(t *testing.T) {
	t.Parallel()

	repo := &bookRepo{books: map[string]book.Book{
		"b-1": {ID: "b-1", Title: "Old", Authors: []string{"Author"}, Status: book.StatusActive},
	}}

	svc := usecase.NewBookService(repo, stubIDGen{id: "ignored"})

	archived, err := svc.SetStatus(context.Background(), "b-1", book.StatusArchived)
	if err != nil {
		t.Fatalf("archive book: %v", err)
	}
	if archived.Status != book.StatusArchived {
		t.Fatalf("expected archived status got %s", archived.Status)
	}

	reactivated, err := svc.SetStatus(context.Background(), "b-1", book.StatusActive)
	if err != nil {
		t.Fatalf("reactivate book: %v", err)
	}
	if reactivated.Status != book.StatusActive {
		t.Fatalf("expected active status got %s", reactivated.Status)
	}
}

func TestMemberServiceRegisterDuplicateID(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	repo := &memberRepo{members: map[string]member.Member{
		"m-1": {ID: "m-1", Name: "Existing", JoinedAt: now, Status: member.StatusActive},
	}}

	svc := usecase.NewMemberService(repo, stubIDGen{id: "m-1"}, stubClock{now: now})

	_, err := svc.Register(context.Background(), dto.RegisterMemberInput{Name: "Joe"})
	if !errors.Is(err, shared.ErrDuplicateID) {
		t.Fatalf("expected duplicate id error got %v", err)
	}
}

func TestMemberServiceUpdate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	repo := &memberRepo{members: map[string]member.Member{
		"m-1": {ID: "m-1", Name: "Old", JoinedAt: now, Status: member.StatusActive},
	}}

	svc := usecase.NewMemberService(repo, stubIDGen{id: "ignored"}, stubClock{now: now})

	updated, err := svc.Update(context.Background(), dto.UpdateMemberInput{ID: "m-1", Name: "New", Email: "new@example.com", Phone: "123"})
	if err != nil {
		t.Fatalf("update member: %v", err)
	}

	if updated.Name != "New" || updated.Email != "new@example.com" || updated.Phone != "123" {
		t.Fatalf("unexpected updated member: %+v", updated)
	}
}

func TestCopyServiceRejectsDuplicateBarcode(t *testing.T) {
	t.Parallel()

	repo := &copyRepo{copies: map[string]copy.Copy{
		"c-1": {ID: "c-1", BookID: "b-1", Barcode: "BC-1", Status: copy.StatusAvailable},
	}}

	svc := usecase.NewCopyService(repo, stubIDGen{id: "c-2"})

	_, err := svc.Create(context.Background(), dto.CreateCopyInput{BookID: "b-1", Barcode: "BC-1"})
	if !errors.Is(err, shared.ErrDuplicateBarcode) {
		t.Fatalf("expected duplicate barcode error got %v", err)
	}
}

func TestCopyServiceUpdate(t *testing.T) {
	t.Parallel()

	repo := &copyRepo{copies: map[string]copy.Copy{
		"c-1": {ID: "c-1", BookID: "b-1", Barcode: "BC-1", Status: copy.StatusAvailable},
	}}

	svc := usecase.NewCopyService(repo, stubIDGen{id: "ignored"})

	updated, err := svc.Update(context.Background(), dto.UpdateCopyInput{ID: "c-1", Barcode: "BC-2", Status: "damaged", ConditionNote: "torn pages"})
	if err != nil {
		t.Fatalf("update copy: %v", err)
	}

	if updated.Status != copy.StatusDamaged || updated.Barcode != "BC-2" {
		t.Fatalf("unexpected updated copy: %+v", updated)
	}
}
