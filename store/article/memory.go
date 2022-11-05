package store

import (
	"context"
	"fmt"

	"github.com/hsmtkk/urban-guacamole/record"
)

type memoryImpl struct {
	internal map[string]record.Article
}

func MemoryImpl() ArticleStore {
	internal := map[string]record.Article{}
	return &memoryImpl{internal}
}

func (m *memoryImpl) Get(ctx context.Context, articleID string) (record.Article, error) {
	article, ok := m.internal[articleID]
	if ok {
		return article, nil
	} else {
		return article, fmt.Errorf("id %s was not found", articleID)
	}
}

func (m *memoryImpl) GetByUserID(ctx context.Context, userID string) ([]record.Article, error) {
	articles := []record.Article{}
	for _, article := range m.internal {
		if userID == article.UserID {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (m *memoryImpl) Set(ctx context.Context, article record.Article) error {
	m.internal[article.ArticleID] = article
	return nil
}

func (m *memoryImpl) Delete(ctx context.Context, articleID string) error {
	delete(m.internal, articleID)
	return nil
}
