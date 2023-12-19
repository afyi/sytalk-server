package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Html      string             `bson:"atContentHtml" validate:"required"`
	Md        string             `bson:"atContentMd" validate:"required"`
	Ua        string             `bson:"ua"`
	Os        string             `bson:"os"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
