package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sreet       string             `bson:"street" json:"street"`
	ShipToName  string             `bson:"ship_to_name" json:"ship_to_name"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Longitude   float64            `bson:"longitude" json:"longitude"`
	Latitude    float64            `bson:"latitude" json:"latitude"`
	Notes       string             `bson:"notes" json:"notes"`
	IsDefault   bool               `bson:"is_default" json:"is_default"`
	RegionID    primitive.ObjectID `bson:"region_id" json:"region_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`

	// Non-persisted field (join-result)
	RegionData *Region `bson:"-" json:"region,omitempty"` // hasil join
}
