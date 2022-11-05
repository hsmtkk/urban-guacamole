package store

import (
	"context"

	"github.com/hsmtkk/urban-guacamole/record"
)

type ArticleStore interface {
	Get(ctx context.Context, articleID string) (record.Article, error)
	GetByUserID(ctx context.Context, userID string) ([]record.Article, error)
	Set(ctx context.Context, article record.Article) error
	Delete(ctx context.Context, articleID string) error
}
