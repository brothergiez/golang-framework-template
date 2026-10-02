package memory

import (
	"context"
	"sync"

	domainuser "github.com/lion/go-framework/internal/domain/user"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
)

type UserRepository struct {
	mu    sync.RWMutex
	byID  map[string]*domainuser.User
	email map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:  make(map[string]*domainuser.User),
		email: make(map[string]string),
	}
}

func (r *UserRepository) Create(_ context.Context, user *domainuser.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.email[user.Email]; ok {
		return apperrors.Conflict("email already exists")
	}
	cp := *user
	r.byID[user.ID] = &cp
	r.email[user.Email] = user.ID
	return nil
}

func (r *UserRepository) FindByID(_ context.Context, id string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.email[email]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	u := r.byID[id]
	cp := *u
	return &cp, nil
}

func (r *UserRepository) List(_ context.Context, offset, limit int) ([]*domainuser.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*domainuser.User, 0, len(r.byID))
	for _, u := range r.byID {
		cp := *u
		all = append(all, &cp)
	}
	total := int64(len(all))
	if offset >= len(all) {
		return []*domainuser.User{}, total, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], total, nil
}

func (r *UserRepository) Update(_ context.Context, user *domainuser.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byID[user.ID]
	if !ok {
		return apperrors.ErrNotFound
	}
	if existing.Email != user.Email {
		delete(r.email, existing.Email)
		r.email[user.Email] = user.ID
	}
	cp := *user
	r.byID[user.ID] = &cp
	return nil
}

func (r *UserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return apperrors.ErrNotFound
	}
	delete(r.email, u.Email)
	delete(r.byID, id)
	return nil
}
