package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	StatusPending   = "pending"
	StatusProcessed = "processed"
)

type OutboxEvent struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Action      string        `bson:"action" json:"action"`
	ActorID     string        `bson:"actorId" json:"actorId"`
	EntityID    string        `bson:"entityId" json:"entityId"`
	Before      any           `bson:"before,omitempty" json:"before,omitempty"`
	After       any           `bson:"after,omitempty" json:"after,omitempty"`
	Status      string        `bson:"status" json:"status"`
	Attempts    int           `bson:"attempts" json:"attempts"`
	LastError   string        `bson:"lastError,omitempty" json:"lastError,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	ProcessedAt *time.Time    `bson:"processedAt,omitempty" json:"processedAt,omitempty"`
}
