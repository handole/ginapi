package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Region struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	State        string             `bson:"state" json:"state"`
	City         string             `bson:"city" json:"city"`
	District     string             `bson:"district" json:"district"`
	Code         string             `bson:"code" json:"code"`
	Zipcode      string             `bson:"zipcode" json:"zipcode"`
	SubDisctrict string             `bson:"sub_district" json:"sub_district"`
}
