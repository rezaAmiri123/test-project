package adapters

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

type MemoryUserRepository struct {
	users map[string]user.User
	mu    sync.RWMutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: map[string]user.User{},
		mu:    sync.RWMutex{},
	}
}

func (m *MemoryUserRepository) Create(ctx context.Context, user *user.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.Username]=*user
	return nil
}

func (m *MemoryUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.users[username]
	if ok {
		return &u, nil
	}
	return nil, errors.New("user not found")
}
