package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Role     string             `json:"role"`
}
