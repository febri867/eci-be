package repository

import (
	"context"
	"time"

	historyModel "eci-backend/internal/history/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OutboxRepository struct {
	collection *mongo.Collection
}

func NewOutboxRepository(db *mongo.Database) *OutboxRepository {
	return &OutboxRepository{
		collection: db.Collection("outbox_events"),
	}
}

func (r *OutboxRepository) Create(ctx context.Context, event historyModel.OutboxEvent) error {
	event.ID = bson.NewObjectID()
	event.Status = historyModel.StatusPending
	event.Attempts = 0
	event.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, event)
	return err
}

func (r *OutboxRepository) FindPending(ctx context.Context, limit int64) ([]historyModel.OutboxEvent, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"status": historyModel.StatusPending},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []historyModel.OutboxEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *OutboxRepository) MarkProcessed(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status":      historyModel.StatusProcessed,
			"processedAt": now,
			"lastError":   "",
		}},
	)

	return err
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, id bson.ObjectID, message string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"attempts": 1},
			"$set": bson.M{"lastError": message},
		},
	)

	return err
}

func (r *OutboxRepository) Record(
	ctx context.Context,
	event historyModel.OutboxEvent,
) error {
	return r.Create(ctx, event)
}
