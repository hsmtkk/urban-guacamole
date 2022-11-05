package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/hsmtkk/urban-guacamole/record"
	"google.golang.org/api/iterator"
)

type firestoreImpl struct {
	projectID  string
	collection string
}

func FirestoreImpl() ArticleStore {
	return &firestoreImpl{}
}

func (f *firestoreImpl) Get(ctx context.Context, articleID string) (record.Article, error) {
	var article record.Article
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return article, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	snap, err := client.Collection(f.collection).Doc(articleID).Get(ctx)
	if err != nil {
		return article, fmt.Errorf("id %s was not found", articleID)
	}
	if err := snap.DataTo(&article); err != nil {
		return article, fmt.Errorf("DocumentSnapshot.DataTo failed; %w", err)
	}
	return article, nil
}

func (f *firestoreImpl) GetByUserID(ctx context.Context, userID string) ([]record.Article, error) {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	articles := []record.Article{}
	iter := client.Collection(f.collection).Where("userid", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("firestore iteration failed; %w", err)
		}
		var article record.Article
		if err := doc.DataTo(&article); err != nil {
			return nil, fmt.Errorf("DocumentSnapshot.DataTo failed; %w", err)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (f *firestoreImpl) Set(ctx context.Context, article record.Article) error {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	if _, err := client.Collection(f.collection).Doc(article.ArticleID).Set(ctx, article); err != nil {
		return fmt.Errorf("DocumentRef.Set failed; %w", err)
	}
	return nil
}

func (f *firestoreImpl) Delete(ctx context.Context, articleID string) error {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	if _, err := client.Collection(f.collection).Doc(articleID).Delete(ctx); err != nil {
		return fmt.Errorf("DocumentRef.Delete failed; %w", err)
	}
	return nil
}
