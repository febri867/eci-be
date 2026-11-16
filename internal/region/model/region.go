package model

import "time"

type Region struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	ParentID  string    `bson:"parentId" json:"parentId"`
	Level     int       `bson:"level" json:"level"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type RegionResponse struct {
	Data []RegionItem `json:"data"`
	Meta RegionMeta   `json:"meta"`
}

type RegionItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RegionMeta struct {
	Level int `json:"level"`
}
