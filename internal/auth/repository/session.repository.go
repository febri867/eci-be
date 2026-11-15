package repository

import (
	"context"
	"time"

	"eci-backend/internal/auth/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection("sessions"),
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *model.Session,
) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *SessionRepository) DeleteByTokenHash(
	ctx context.Context,
	tokenHash string,
) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"token_hash": tokenHash,
	})
	return err
}

func (r *SessionRepository) DeleteExpired(
	ctx context.Context,
	now time.Time,
) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lte": now},
	})
	return err
}

func (r *SessionRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*model.Session, error) {
	var session model.Session

	err := r.collection.FindOne(ctx, bson.M{
		"token_hash": tokenHash,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&session)

	if err != nil {
		return nil, err
	}

	return &session, nil
}
