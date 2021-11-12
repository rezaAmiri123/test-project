package article

import "context"

type Repository interface {
	Create(ctx context.Context, article *Article) error
	GetBySlug(ctx context.Context, slug string)(*Article,error)
}