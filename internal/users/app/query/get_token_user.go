package query

import (
	"context"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

type GetUserTokenHandler struct {
	repo user.Repository
}

func NewGetUserTokenHandler(userRepo user.Repository) GetUserTokenHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}
	return GetUserTokenHandler{repo: userRepo}
}

func (h GetUserTokenHandler) Handler(ctx context.Context, token string) (*user.User, error) {
	username, err := user.GetUsernameFromJWTToken(token)
	if err != nil{
		return &user.User{}, err
	}
	u,err := h.repo.GetByUsername(ctx,username)
	if err != nil{
		return &user.User{}, err
	}
	u.HidePassword()
	return u,nil
}