package repository

import (
	"context"

	"eci-backend/internal/auth/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var user model.User

	err := r.collection.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.User, error) {
	var user model.User

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
