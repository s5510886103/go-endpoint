package main

import (
	"fmt"
	_ "github.com/hlkittipan/go-endpoint/docs"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/controller"
	"github.com/hlkittipan/go-endpoint/src/middleware"
	"github.com/hlkittipan/go-endpoint/src/service"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"                  // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
	"log"
	"net/http"
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
	var loginService = service.StaticLoginService()
	var jwtService = service.JWTAuthService()
	var loginController = controller.LoginHandler(loginService, jwtService)

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

	r.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	h := CustomerHandler{}
	h.Initialize()

	v1 := r.Group("/v1")
	v1.Use(middleware.AuthorizeJWT())
	{
		v1.GET("/customers", h.GetAllCustomer)
		v1.GET("/customers/:id", h.GetCustomer)
		v1.POST("/customers", h.SaveCustomer)
		v1.PUT("/customers/:id", h.UpdateCustomer)
		v1.DELETE("/customers/:id", h.DeleteCustomer)
	}

	r.GET("/users", controller.GetAllUsers())           //add this
	r.POST("/user", controller.CreateUser())            //add this
	r.GET("/user/:userId", controller.GetAUser())       //add this
	r.PUT("/user/:userId", controller.EditAUser())      //add this
	r.DELETE("/user/:userId", controller.DeleteAUser()) //add this

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

type CustomerHandler struct {
	DB *gorm.DB
}

type Customer struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}

func (h *CustomerHandler) Initialize() {
	//dbUsername := goDotEnvVariable("MYSQL_USERNAME")
	//dbPassword := goDotEnvVariable("MYSQL_PASSWORD")
	//dbName := goDotEnvVariable("MYSQL_DB_NAME")
	//dbHost := goDotEnvVariable("MYSQL_HOST")
	//dbPort := goDotEnvVariable("MYSQL_PORT")
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
	//fmt.Println(dsn)
	//db, err := gorm.Open(mysql.Open("root:1234567890@tcp(127.0.0.1:3306)/golang?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = db.AutoMigrate(&Customer{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//h.DB = db
}

// GetAllCustomer godoc
// @summary Get customer list
// @description Get all customer
// @tags customers
// @id CustomerHandler
// @accept json
// @produce json
// @response 200 {object} model.Response "OK"
// @response 400 {object} model.Response "Bad Request"
// @response 401 {object} model.Response "Unauthorized"
// @response 500 {object} model.Response "Internal Server Error"
// @router /customers [get]
func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	customers := []Customer{}

	h.DB.Find(&customers)
	fmt.Printf("ClientIP: %s\n", c.ClientIP())
	// If the client is 192.168.1.2, use the X-Forwarded-For
	// header to deduce the original client IP from the trust-
	// worthy parts of that header.
	// Otherwise, simply return the direct client IP
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer godoc
// @summary Update Customer
// @description Update customer by id
// @tags customers
// @security ApiKeyAuth
// @id UpdateCustomer
// @accept json
// @produce json
// @param id path int true "id of customer to be updated"
// @param Customer body model.CustomerForUpdate true "Customer data to be updated"
// @response 200 {object} model.Response "OK"
// @response 400 {object} model.Response "Bad Request"
// @response 401 {object} model.Response "Unauthorized"
// @response 500 {object} model.Response "Internal Server Error"
// @Router /customers/:id [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
