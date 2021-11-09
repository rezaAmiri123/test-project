package auth

import (
	"context"
	"github.com/pkg/errors"
)



type User struct {
	Username  string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	// if we expect that the user of the function may be interested with concrete error,
	// it's a good idea to provide variable with this error
	NoUserInContextError = errors.New("no user in context")
)

func UserFromContext(ctx context.Context) (User, error) {
	u,ok:=ctx.Value(userContextKey).(User)
	if ok{
		return u,NoUserInContextError
	}
	return User{}, NoUserInContextError
}
