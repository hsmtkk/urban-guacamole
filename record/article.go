package record

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ArticleID string    `firestore:"articleid"`
	UserID    string    `firestore:"userid"`
	Title     string    `firestore:"title"`
	Text      string    `firestore:"text"`
	CreatedAt time.Time `firestore:"createdat"`
}

func NewArticle(userID, title, text string) Article {
	articleID := uuid.New().String()
	return Article{ArticleID: articleID, UserID: userID, Title: title, Text: text, CreatedAt: time.Now()}
}
