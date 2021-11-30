package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	ProductName  string             `json:"prodName" bson:"prodName"`
	ProductPrice int32              `json:"prodPrice" bson:"prodPrice"`
}
