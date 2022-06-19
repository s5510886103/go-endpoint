package main

import (
	"fmt"
	_ "github.com/hlkittipan/go-endpoint/docs"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/controller"
	"github.com/hlkittipan/go-endpoint/src/middleware"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"                  // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"os"
)

import "github.com/gin-gonic/gin"

// @title Customers API
// @version 1.0
// @description start
// @termsOfService not yet

// @contact.name API Support
// @contact.url not yet
// @contact.email hlkittipan@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes https http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fmt.Println("Hello world")
	fmt.Println("KS.")
	//run database
	config.ConnectDB()
	r := setupRouter()
	err := r.Run(":5555")
	if err != nil {
		fmt.Println("Error starter")
		return
	} // listen and serve on 0.0.0.0:5555 (for windows "localhost:5555")

}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	//r.SetTrustedProxies([]string{"IP SERVER"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", controller.Login())

	v1 := r.Group("/v1")
	v1.Use(middleware.Authentication())
	{
		v1.GET("/customer/:id", controller.GetACustomer())
		v1.POST("/customer", controller.CreateCustomer())
		v1.PUT("/customer/:id", controller.EditACustomer())
		v1.DELETE("/customer/:id", controller.DeleteACustomer())

		v1.DELETE("/user/:userId", controller.DeleteAUser())
	}
	r.GET("/customers", controller.GetAllCustomers())

	r.GET("/users", controller.GetAllUsers())
	r.POST("/user", controller.CreateUser())       //add this
	r.GET("/user/:userId", controller.GetAUser())  //add this
	r.PUT("/user/:userId", controller.EditAUser()) //add this

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
