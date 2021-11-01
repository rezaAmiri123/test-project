package query

import (
	"context"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

type GetProfileHandler struct {
	userRepo user.Repository
}

func NewGetProfileHandler(userRepo user.Repository) GetProfileHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}
	return GetProfileHandler{userRepo: userRepo}
}

func (h GetProfileHandler) Handler(ctx context.Context, username string) (*user.User, error) {
	u, err := h.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return &user.User{}, err
	}
	return u, nil
}
