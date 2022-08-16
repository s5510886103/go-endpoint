package main

import (
	"fmt"
	_ "github.com/hlkittipan/go-endpoint/docs"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/controller"
	"github.com/hlkittipan/go-endpoint/src/middleware"
	"github.com/swaggo/files"                  // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"os"
	"os/signal"
	"syscall"
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
	log.Print("Starting the service...")

	gin.SetMode(config.GoDotEnvVariable("GIN_MODE"))
	//run database
	config.ConnectDB()
	r := setupRouter()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	log.Print("The service is ready to listen and serve.")

	err := r.Run(":" + config.GoDotEnvVariable("PORT"))
	if err != nil {
		fmt.Println("Error starter")
		return
	} // listen and serve on 0.0.0.0:5555 (for windows "localhost:5555")
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	log.Print("Done")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
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

	r.GET("/home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello! Your request was processed.",
		})
	},
	)

	r.POST("/login", controller.Login())

	v1 := r.Group("/v1")
	v1.Use(middleware.Authentication())
	{
		v1.GET("/customer/:id", controller.GetACustomer())
		v1.POST("/customer", controller.CreateCustomer())
		v1.PUT("/customer/:id", controller.EditACustomer())
		v1.DELETE("/customer/:id", controller.DeleteACustomer())
		v1.GET("/customers", controller.GetAllCustomers())

		v1.DELETE("/user/:userId", controller.DeleteAUser())
		v1.GET("/users", controller.GetAllUsers())
		v1.GET("/user/:userId", controller.GetAUser())  //add this
		v1.PUT("/user/:userId", controller.EditAUser()) //add this

	}

	r.POST("/user", controller.CreateUser()) //add this
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
