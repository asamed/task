package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductName  string             `json:"prodName" bson:"prodName"`
	ProductPrice int                `json:"prodPrice" bson:"prodPrice"`
}
