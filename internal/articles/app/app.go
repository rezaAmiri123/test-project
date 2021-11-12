package app

import (
	"github.com/rezaAmiri123/test-project/internal/articles/app/command"
	"github.com/rezaAmiri123/test-project/internal/articles/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Queries struct {
	GetArticle query.GetArticleHandler
}

type Commands struct {
	CreateArticle command.CreateArticleHandler
}
