package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//We can see here that we used * pointer for some and others are not this is beacause
// * can hold nil values so === simply (optional)
// without * they represent required values

type Food struct {
	ID         primitive.ObjectID `bson:"id"`
	Name       *string            `json:"name" validate:"required,min=2,max=100"`
	Price      *float64           `json:"price" validate:"required"`
	Food_image *string            `json:"food_image" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Food_id    string             `json:"food_id"`
	Menu_id    *string            `json:"menu_id" validate:"required"`
}
