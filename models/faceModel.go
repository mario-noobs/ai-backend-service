package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Face struct {
	ID        *primitive.ObjectID `bson:"_id"`
	Name      *string             `json:"name" validate:"required,min=2,max=100"`
	Gender    *float64            `json:"gender"`
	Image     *string             `json:"image" validate:"required"`
	CreatedAt *time.Time          `json:"created_at"`
	UpdatedAt *time.Time          `json:"updated_at"`
}
