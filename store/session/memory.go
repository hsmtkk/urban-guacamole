package session

import (
	"context"
	"fmt"

	"github.com/hsmtkk/urban-guacamole/record"
)

type memoryImpl struct {
	internal map[string]record.Session
}

func MemoryImpl() SessionStore {
	internal := map[string]record.Session{}
	return &memoryImpl{internal}
}

func (m *memoryImpl) Get(ctx context.Context, sessionID string) (record.Session, error) {
	session, ok := m.internal[sessionID]
	if ok {
		return session, nil
	} else {
		return session, fmt.Errorf("id %s was not found", sessionID)
	}
}

func (m *memoryImpl) Set(ctx context.Context, session record.Session) error {
	m.internal[session.SessionID] = session
	return nil
}

func (m *memoryImpl) Delete(ctx context.Context, sessionID string) error {
	delete(m.internal, sessionID)
	return nil
}
