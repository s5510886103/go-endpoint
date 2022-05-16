package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CustomerResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Customer struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Firstname string             `json:"firstname" binding:"required" example:"Choopong" maxLength:"255"`   // Customer Firstname
	Lastname  string             `json:"lastname" binding:"required" example:"Choosamer" maxLength:"255"`   // Customer Lastname
	Email     string             `json:"email" binding:"required" example:"choo@gmail.com" maxLength:"255"` // Customer E-mail
	Gender    string             `json:"gender" example:"male" enums:"male,female"`                         // Customer Gender
}
