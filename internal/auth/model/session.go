package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Session struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id"`
	TokenHash string        `bson:"token_hash"`
	ExpiresAt time.Time     `bson:"expires_at"`
	CreatedAt time.Time     `bson:"created_at"`
}
