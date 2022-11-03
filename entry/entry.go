package entry

import "time"

type Entry struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdat"`
}

func New(title, text string) Entry {
	id := time.Now().Unix()
	return Entry{
		ID:        id,
		Title:     title,
		Text:      text,
		CreatedAt: time.Now(),
	}
}
