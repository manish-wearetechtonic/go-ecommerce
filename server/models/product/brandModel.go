package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Brand struct {
	ID         primitive.ObjectID `bson:"_id"`
	Brand_name string             `json:"brand_name" validate:"required"`
}
