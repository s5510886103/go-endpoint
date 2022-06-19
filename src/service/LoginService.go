package service

import (
	"context"
	"fmt"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

	passwordIsValid, _ := VerifyPassword(*foundUser.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		return nil
	}

	return &loginInformation{
		email:    *foundUser.Email,
		password: *foundUser.Password,
	}
}

//VerifyPassword checks the input password while verifying it with the password in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func (info *loginInformation) LoginUser(email string, password string) bool {
	fmt.Println(info)
	fmt.Println(password)
	return info.email == email && info.password == password
}
