package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Employee struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Position  string        `bson:"position" json:"position"`
	Grade     string        `bson:"grade" json:"grade"`
	Salary    int64         `bson:"salary" json:"salary"`
	Address   string        `bson:"address" json:"address"`
	Province  string        `bson:"province" json:"province"`
	City      string        `bson:"city" json:"city"`
	District  string        `bson:"district" json:"district"`
	BirthDate time.Time     `bson:"birthDate" json:"birthDate"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
