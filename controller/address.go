package controller

import (
	"context"
	"ecomm/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, "invalid search index")
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "server error")
		}
		var addresses model.Address
		addresses.Address_id = primitive.NewObjectID()
		if err = c.Bind(&addresses); err != nil {
			return c.JSON(http.StatusBadGateway, err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			return c.JSON(500, err.Error())
		}
		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			return c.JSON(500, err.Error())
		}
		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: "addresses"}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return c.JSON(500, "could not update")
			}
		} else {
			return c.JSON(400, "Not allowed")
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "updated successfully")
	}
}
func EditHomeAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, "invalid search index")
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "server err")
		}
		var editAddress model.Address
		if err = c.Bind(&editAddress); err != nil {
			return c.JSON(http.StatusBadRequest, "invalid request")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(500, "could not edit address")
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "Edited")
	}
}
func EditWorkAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, "invalid search index")
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "server err")
		}
		var editAddress model.Address
		if err = c.Bind(&editAddress); err != nil {
			return c.JSON(http.StatusBadRequest, "invalid request")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(500, "something Went wrong")

		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "Successfully updated the Work Address")

	}
}
func DeleteAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, "invalid search index")
		}
		addresses := make([]model.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "server err")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(404, "Wrong command")
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "deleted")
	}
}
