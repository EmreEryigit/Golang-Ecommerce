package controller

import (
	"context"
	"ecomm/database"
	"ecomm/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection *mongo.Collection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		productQueryId := c.QueryParam("id")
		if productQueryId == "" {
			return c.JSON(http.StatusBadRequest, "product id is empty")
		}
		userQueryId := c.QueryParam("userID")
		if userQueryId == "" {
			return c.JSON(http.StatusBadRequest, "user id is empty")
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(200, "Added to cart")
	}
}

func (app *Application) RemoveItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		productQueryId := c.QueryParam("id")
		if productQueryId == "" {
			return c.JSON(http.StatusBadRequest, "product id is empty")
		}
		userQueryId := c.QueryParam("userID")
		if userQueryId == "" {
			return c.JSON(http.StatusBadRequest, "user id is empty")
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(200, "Added to cart")
	}
}

func GetItemFromCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		userQueryId := c.QueryParam("userID")
		if userQueryId == "" {
			return c.JSON(http.StatusBadRequest, "user id is empty")
		}
		usert_id, _ := primitive.ObjectIDFromHex(userQueryId)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var filledCart model.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
		if err != nil {
			return c.JSON(404, "Not found")
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, group})
		if err != nil {
			return c.JSON(500, "server error")
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			return c.JSON(500, "server error")
		}

		for _, json := range listing {
			c.JSONPretty(200, filledCart.UserCart, "4")
			c.JSON(200, json["total"])
		}
		ctx.Done()
		return err
	}
}

func (app *Application) BuyFromCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		userQueryId := c.QueryParam("userID")
		if userQueryId == "" {
			return c.JSON(http.StatusBadRequest, "user id is empty")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(200, "placed the order")
	}
}

func (app *Application) InstantBuy() echo.HandlerFunc {
	return func(c echo.Context) error {
		productQueryId := c.QueryParam("id")
		if productQueryId == "" {
			return c.JSON(http.StatusBadRequest, "product id is empty")
		}
		userQueryId := c.QueryParam("userID")
		if userQueryId == "" {
			return c.JSON(http.StatusBadRequest, "user id is empty")
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(200, "removed from cart")
	}
}
