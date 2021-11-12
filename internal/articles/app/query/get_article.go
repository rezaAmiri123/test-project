package query

import (
	"context"
	"github.com/rezaAmiri123/test-project/internal/articles/domain/article"
)

type GetArticleHandler struct {
	articleRepo article.Repository
}

func NewGetArticleHandler(articleRepo article.Repository) GetArticleHandler {
	if articleRepo == nil {
		panic("article repo is nill")
	}
	return GetArticleHandler{articleRepo: articleRepo}
}

func (h GetArticleHandler) Handle(ctx context.Context, slug string) (*article.Article, error) {
	a,err:= h.articleRepo.GetBySlug(ctx, slug)
	if err != nil{
		return &article.Article{}, err
	}
	return a,nil
}
