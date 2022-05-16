package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hlkittipan/go-endpoint/src/config"
	"github.com/hlkittipan/go-endpoint/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var customerCollection *mongo.Collection = config.GetCollection(config.DB, "customers")

func CreateCustomer() gin.HandlerFunc {
	var validate = validator.New()

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var customer model.Customer
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, model.CustomerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&customer); validationErr != nil {
			c.JSON(http.StatusBadRequest, model.CustomerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newCustomer := model.Customer{
			Id:        primitive.NewObjectID(),
			Firstname: customer.Firstname,
			Lastname:  customer.Lastname,
			Email:     customer.Email,
			Gender:    customer.Gender,
		}

		result, err := customerCollection.InsertOne(ctx, newCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, model.CustomerResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetACustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		customerId := c.Param("customerId")
		var customer model.Customer
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(customerId)

		err := customerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, model.CustomerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": customer}})
	}
}

// EditACustomer godoc
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
func EditACustomer() gin.HandlerFunc {
	var validate = validator.New()
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		customerId := c.Param("customerId")
		var customer model.Customer
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(customerId)

		//validate the request body
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, model.CustomerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&customer); validationErr != nil {
			c.JSON(http.StatusBadRequest, model.CustomerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"firstname": customer.Firstname, "lastname": customer.Lastname, "email": customer.Email, "gender": customer.Gender}
		result, err := customerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated customer details
		var updatedCustomer model.Customer
		if result.MatchedCount == 1 {
			err := customerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedCustomer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, model.CustomerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedCustomer}})
	}
}

func DeleteACustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		customerId := c.Param("customerId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(customerId)

		result, err := customerCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				model.CustomerResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Customer} with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			model.CustomerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Customer} successfully deleted!"}},
		)
	}
}

// GetAllCustomers godoc
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
func GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var customers []model.Customer
		defer cancel()

		results, err := customerCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer func(results *mongo.Cursor, ctx context.Context) {
			err := results.Close(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}(results, ctx)
		for results.Next(ctx) {
			var singleCustomer model.Customer
			if err = results.Decode(&singleCustomer); err != nil {
				c.JSON(http.StatusInternalServerError, model.CustomerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			customers = append(customers, singleCustomer)
		}

		c.JSON(http.StatusOK,
			model.CustomerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": customers}},
		)
	}
}
