package user

import (
	"context"
	"errors"
	"strings"
	"time"

	domainuser "github.com/lion/go-framework/internal/domain/user"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
	"github.com/lion/go-framework/internal/shared/pagination"
	"github.com/lion/go-framework/internal/shared/utils"
)

type Service struct {
	repo domainuser.Repository
}

func NewService(repo domainuser.Repository) *Service {
	return &Service{repo: repo}
}

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	Name  string
	Email string
}

func (s *Service) Create(ctx context.Context, in CreateUserInput) (*domainuser.User, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if name == "" || email == "" {
		return nil, apperrors.Validation("name and email are required", nil)
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, apperrors.Internal("failed to check existing user", err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("email already exists")
	}

	now := time.Now().UTC()
	u := &domainuser.User{
		ID:        utils.NewID(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, apperrors.Internal("failed to create user", err)
	}
	return u, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*domainuser.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NotFound("user not found")
		}
		return nil, apperrors.Internal("failed to get user", err)
	}
	return u, nil
}

func (s *Service) List(ctx context.Context, params pagination.Params) (*pagination.Result[*domainuser.User], error) {
	items, total, err := s.repo.List(ctx, params.Offset(), params.Limit)
	if err != nil {
		return nil, apperrors.Internal("failed to list users", err)
	}
	return &pagination.Result[*domainuser.User]{
		Items: items,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

func (s *Service) Update(ctx context.Context, id string, in UpdateUserInput) (*domainuser.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NotFound("user not found")
		}
		return nil, apperrors.Internal("failed to get user", err)
	}

	if name := strings.TrimSpace(in.Name); name != "" {
		u.Name = name
	}
	if email := strings.TrimSpace(strings.ToLower(in.Email)); email != "" {
		existing, err := s.repo.FindByEmail(ctx, email)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.Internal("failed to check existing user", err)
		}
		if existing != nil && existing.ID != id {
			return nil, apperrors.Conflict("email already exists")
		}
		u.Email = email
	}
	u.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, u); err != nil {
		return nil, apperrors.Internal("failed to update user", err)
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NotFound("user not found")
		}
		return apperrors.Internal("failed to get user", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal("failed to delete user", err)
	}
	return nil
}
