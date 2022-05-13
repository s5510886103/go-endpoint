package main

import (
	"fmt"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
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
	r := setupRouter()
	err := r.Run(":4000")
	if err != nil {
		fmt.Println("Error starter")
		return
	} // listen and serve on 0.0.0.0:4000 (for windows "localhost:4000")

}

func setupRouter() *gin.Engine {
	r := gin.Default()

	err := r.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	//r.SetTrustedProxies([]string{"IP SERVER"})

	h := CustomerHandler{}
	h.Initialize()

	r.GET("/customers", h.GetAllCustomer)
	r.GET("/customers/:id", h.GetCustomer)
	r.POST("/customers", h.SaveCustomer)
	r.PUT("/customers/:id", h.UpdateCustomer)
	r.DELETE("/customers/:id", h.DeleteCustomer)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	dsn := "root:blockee-dev@tcp(127.0.0.1:3306)/golang?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Customer{})
	if err != nil {
		log.Fatal(err)
	}

	h.DB = db
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
