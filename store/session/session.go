package session

import (
	"context"

	"github.com/hsmtkk/urban-guacamole/record"
)

type SessionStore interface {
	Get(ctx context.Context, sessionID string) (record.Session, error)
	Set(ctx context.Context, data record.Session) error
	Delete(ctx context.Context, sessionID string) error
}
