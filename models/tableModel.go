package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID               primitive.ObjectID `bson:"id"`
	Number_of_guests *int               `json:"number_of_guests" valdate:"required"`
	Table_number     *int               `json:"table_number" validate:"required"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
	Table_id         string             `json:"table_id"`
	Status           *string            `json:"status" validate:"eq=OCCUPIED|eq=VACANT"`
}
