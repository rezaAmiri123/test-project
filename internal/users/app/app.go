package app

import (
	"github.com/rezaAmiri123/test-project/internal/users/app/command"
	"github.com/rezaAmiri123/test-project/internal/users/app/query"
	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
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

func NewApplication(repository user.Repository) *Application {
	return &Application{
		Commands: Commands{
			CreateUser: command.NewCreateUserHandler(repository),
		},
		Queries: Queries{
			GetProfile: query.NewGetProfileHandler(repository),
		},
	}
}
