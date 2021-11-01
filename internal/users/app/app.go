package app

import (
	"github.com/rezaAmiri123/test-project/internal/users/app/command"
	"github.com/rezaAmiri123/test-project/internal/users/app/query"
)

type Application struct {
	Commands Commands
	Queries Queries
}

type Commands struct {
	CreateUser command.CreateUserHandler
}

type Queries struct {
	GetProfile query.GetProfileHandler
}
