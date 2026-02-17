package ports

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/domain/book"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/copy"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/loan"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
)

type BookRepository interface {
	Save(ctx context.Context, b book.Book) error
	GetByID(ctx context.Context, id string) (book.Book, error)
}

type CopyRepository interface {
	Save(ctx context.Context, c copy.Copy) error
	GetByID(ctx context.Context, id string) (copy.Copy, error)
}

type MemberRepository interface {
	Save(ctx context.Context, m member.Member) error
	GetByID(ctx context.Context, id string) (member.Member, error)
}

type LoanRepository interface {
	Save(ctx context.Context, l loan.Loan) error
	GetByID(ctx context.Context, id string) (loan.Loan, error)
	CountActiveByMemberID(ctx context.Context, memberID string) (int, error)
}
