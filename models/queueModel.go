package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Queue struct {
	ID            primitive.ObjectID `bson:"_id"`
	Queue_id      string             `json:"queue_id"`
	Current_order int                `json:"current_order"`
	Total_orders  int                `json:"total_orders"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
