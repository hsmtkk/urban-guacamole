package record

import "github.com/google/uuid"

type Session struct {
	SessionID string `firestore:"sessionid"`
	UserID    string `firestore:"userid"`
}

func NewSession(userID string) Session {
	sessionID := uuid.New().String()
	return Session{SessionID: sessionID, UserID: userID}
}
