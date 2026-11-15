package repository

import (
	"context"
	"errors"
	"time"

	"eci-backend/internal/employee/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrEmployeeNotFound = errors.New("employee not found")

type EmployeeRepository struct {
	collection *mongo.Collection
}

func NewEmployeeRepository(db *mongo.Database) *EmployeeRepository {
	return &EmployeeRepository{
		collection: db.Collection("employees"),
	}
}

func (r *EmployeeRepository) Create(
	ctx context.Context,
	employee *model.Employee,
) error {
	now := time.Now()

	employee.ID = bson.NewObjectID()
	employee.CreatedAt = now
	employee.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, employee)
	return err
}

func (r *EmployeeRepository) FindAll(
	ctx context.Context,
	filter bson.M,
	page int64,
	limit int64,
) ([]model.Employee, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	employees := make([]model.Employee, 0)

	if err := cursor.All(ctx, &employees); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *EmployeeRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*model.Employee, error) {
	var employee model.Employee

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&employee)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrEmployeeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepository) Update(
	ctx context.Context,
	id bson.ObjectID,
	employee *model.Employee,
) error {
	employee.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":      employee.Name,
			"position":  employee.Position,
			"grade":     employee.Grade,
			"salary":    employee.Salary,
			"address":   employee.Address,
			"province":  employee.Province,
			"city":      employee.City,
			"district":  employee.District,
			"birthDate": employee.BirthDate,
			"updatedAt": employee.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}

func (r *EmployeeRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {
	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": id},
	)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}
