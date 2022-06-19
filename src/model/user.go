package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type User struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Name         *string            `json:"name" validate:"required,min=2,max=200"`
	Email        *string            `json:"email" validate:"email,required"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserId       string             `json:"user_id"`
}
