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
	var user model.User
	var foundUser model.User

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		return nil
	}

	passwordIsValid, msg := controller.VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		return nil
	}

	return &loginInformation{
		email:    "bikash.dulal@wesionary.team",
		password: "testing",
	}
}
func (info *loginInformation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}
