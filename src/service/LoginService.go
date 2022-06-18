package service

import (
	"context"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/controller"
	"github.com/hlkittipan/go-endpoint/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoginService interface {
	LoginUser(email string, password string) bool
}
type loginInformation struct {
	email    string
	password string
}

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")

func StaticLoginService() LoginService {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var foundUser model.User

	err := userCollection.FindOne(ctx, bson.M{"email": foundUser.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return nil
	}

	passwordIsValid, _ := controller.VerifyPassword(*foundUser.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		return nil
	}

	return &loginInformation{
		email:    *foundUser.Email,
		password: *foundUser.Password,
	}
}
func (info *loginInformation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}
