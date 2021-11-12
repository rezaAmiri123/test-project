package command

import (
	"context"
	"github.com/gosimple/slug"
	"github.com/google/uuid"
	"github.com/rezaAmiri123/test-project/internal/articles/domain/article"
)

type CreateArticleHandler struct {
	articleRepo article.Repository
}

func NewCreateArticleHandler(articleRepo article.Repository) CreateArticleHandler {
	if articleRepo == nil {
		panic("articleRepo is nil")
	}
	return CreateArticleHandler{articleRepo: articleRepo}
}

func (h CreateArticleHandler) Handle(ctx context.Context, article *article.Article, userUUID string) error {
	if err := article.Validate(); err != nil {
		return err
	}
	article.UUID = uuid.New().String()
	article.UserUUID = userUUID
	article.Slug =slug.Make(article.Title)

	err := h.articleRepo.Create(ctx, article)
	return err
}
