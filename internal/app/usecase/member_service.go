package usecase

import (
	"context"

	"github.com/mibienpanjoe/LMS-bit/internal/app/dto"
	"github.com/mibienpanjoe/LMS-bit/internal/app/ports"
	"github.com/mibienpanjoe/LMS-bit/internal/domain/member"
)

type MemberService struct {
	members ports.MemberRepository
	idGen   ports.IDGenerator
	clock   ports.Clock
}

func NewMemberService(members ports.MemberRepository, idGen ports.IDGenerator, clock ports.Clock) MemberService {
	return MemberService{members: members, idGen: idGen, clock: clock}
}

func (s MemberService) Register(ctx context.Context, input dto.RegisterMemberInput) (member.Member, error) {
	id := input.ID
	if id == "" {
		id = s.idGen.NewID()
	}

	m := member.Member{
		ID:       id,
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		JoinedAt: s.clock.Now(),
		Status:   member.StatusActive,
	}

	if err := m.Validate(); err != nil {
		return member.Member{}, err
	}

	if err := s.members.Save(ctx, m); err != nil {
		return member.Member{}, err
	}

	return m, nil
}

func (s MemberService) SetStatus(ctx context.Context, memberID string, status member.Status) (member.Member, error) {
	m, err := s.members.GetByID(ctx, memberID)
	if err != nil {
		return member.Member{}, err
	}

	m.Status = status
	if err := m.Validate(); err != nil {
		return member.Member{}, err
	}

	if err := s.members.Save(ctx, m); err != nil {
		return member.Member{}, err
	}

	return m, nil
}

func (s MemberService) List(ctx context.Context) ([]member.Member, error) {
	return s.members.List(ctx)
}
