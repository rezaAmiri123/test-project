package command

import (
	"context"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
)

type CreateUserHandler struct {
	userRepo user.Repository
}

func NewCreateUserHandler(userRepo user.Repository) CreateUserHandler {
	if userRepo==nil{
		panic("userRepo is nil")
	}
	return CreateUserHandler{userRepo: userRepo}
}

func (h CreateUserHandler) Handle(ctx context.Context, user *user.User)error  {
	if err:=user.SetUUID();err!=nil{
		return err
	}
	if err :=user.Validate();err!= nil{
		return err
	}
	if err:=user.HashPassword();err!= nil{
		return err
	}
	err := h.userRepo.Create(ctx, user)
	return err
}
