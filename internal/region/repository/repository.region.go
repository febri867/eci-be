package repository

import (
	"context"

	"eci-backend/internal/region/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RegionRepository struct {
	collection *mongo.Collection
}

func NewRegionRepository(db *mongo.Database) *RegionRepository {
	return &RegionRepository{
		collection: db.Collection("regions"),
	}
}

func (r *RegionRepository) FindByLevelAndParent(
	ctx context.Context,
	level int,
	parentID string,
) ([]model.Region, error) {
	filter := bson.M{
		"level": level,
	}

	if parentID != "" {
		filter["parentId"] = parentID
	}

	findOptions := options.Find().SetSort(bson.D{
		{Key: "name", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	regions := make([]model.Region, 0)

	if err := cursor.All(ctx, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func (r *RegionRepository) UpsertMany(
	ctx context.Context,
	regions []model.Region,
) error {
	if len(regions) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(regions))

	for _, region := range regions {
		updateModel := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": region.ID}).
			SetUpdate(bson.M{
				"$set": bson.M{
					"name":      region.Name,
					"parentId":  region.ParentID,
					"level":     region.Level,
					"updatedAt": region.UpdatedAt,
				},
				"$setOnInsert": bson.M{
					"createdAt": region.CreatedAt,
				},
			}).
			SetUpsert(true)

		models = append(models, updateModel)
	}

	_, err := r.collection.BulkWrite(ctx, models)
	return err
}
