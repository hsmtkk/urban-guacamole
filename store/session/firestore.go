package session

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/hsmtkk/urban-guacamole/record"
)

type firestoreImpl struct {
	projectID  string
	collection string
}

func FirestoreImpl() SessionStore {
	return &firestoreImpl{}
}

func (f *firestoreImpl) Get(ctx context.Context, sessionID string) (record.Session, error) {
	var session record.Session
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return session, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	snap, err := client.Collection(f.collection).Doc(sessionID).Get(ctx)
	if err != nil {
		return session, fmt.Errorf("id %s was not found", sessionID)
	}
	if err := snap.DataTo(&session); err != nil {
		return session, fmt.Errorf("DocumentSnapshot.DataTo failed; %w", err)
	}
	return session, nil
}

func (f *firestoreImpl) Set(ctx context.Context, session record.Session) error {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	if _, err := client.Collection(f.collection).Doc(session.SessionID).Set(ctx, session); err != nil {
		return fmt.Errorf("DocumentRef.Set failed; %w", err)
	}
	return nil
}

func (f *firestoreImpl) Delete(ctx context.Context, sessionID string) error {
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	if _, err := client.Collection(f.collection).Doc(sessionID).Delete(ctx); err != nil {
		return fmt.Errorf("DocumentRef.Delete failed; %w", err)
	}
	return nil
}
