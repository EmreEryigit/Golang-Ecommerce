package controller

import (
	"context"
	"ecomm/database"
	"ecomm/helper"
	"ecomm/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "Users")
var ProductCollection *mongo.Collection = database.OpenCollection(database.Client, "Products")

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		// first initialize private user for password validation
		var userPrivate model.UserPrivate
		if err := c.Bind(&userPrivate); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid request")
		}
		// validate

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": userPrivate.Email})
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "user does not exist")
		}
		if count > 0 {
			defer cancel()
			return c.JSON(http.StatusConflict, "email already taken")
		}
		userPrivate.HashPassword()
		userPrivate.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userPrivate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userPrivate.ID = primitive.NewObjectID()
		userPrivate.User_ID = userPrivate.ID.Hex()
		userPrivate.UserCart = make([]model.ProductUser, 0)
		userPrivate.Address_Details = make([]model.Address, 0)
		userPrivate.Order_Status = make([]model.Order, 0)
		user := userPrivate.User
		validationError := validate.Struct(user)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "could not save user")
		}

		jwtToken, err := helper.GenerateJWT(userPrivate.User_ID, *userPrivate.First_Name, *userPrivate.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating jwt token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = jwtToken
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			defer cancel()
			c.JSON(http.StatusInternalServerError, "error while generating jwt token")
			return err
		}
		defer cancel()
		return c.JSON(http.StatusOK, user)
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var user model.UserPrivate
		var foundUser model.User
		if err := c.Bind(&user); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "user does not exist")
		}
		if foundUser.Email == nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "user not found")
		}
		isValid := foundUser.VerifyPassword(*user.Password)
		if !isValid {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid email or password")
		}
		token, err := helper.GenerateJWT(fmt.Sprint(foundUser.ID), *foundUser.First_Name, *foundUser.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = token
		err1 := session.Save(c.Request(), c.Response())
		if err1 != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "could not save the cookie")
		}
		defer cancel()
		return c.JSON(http.StatusOK, foundUser)
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := Store.Get(c.Request(), "auth-session")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error sessions")
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error saving session")
		}
		return err
	}
}

func ProductViewerAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products model.Product
		defer cancel()
		if err := c.Bind(&products); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())

		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			return c.JSON(http.StatusInternalServerError, "Not Created")

		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
		return nil
	}
}

func SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var productList []model.Product
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.M{})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "please retry")
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "cannot read products")
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			return c.JSON(400, "invalid")
		}
		defer cancel()
		return c.JSON(200, productList)

	}
}

func SearchProductByQuery() echo.HandlerFunc {
	return func(c echo.Context) error {
		var searchProducts []model.Product
		queryParam := c.QueryParam("name")
		if queryParam == "" {
			log.Println("query is empty")
			return c.JSON(http.StatusNotFound, "invalid search index")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "could not fetch products")
		}
		err = cursor.All(ctx, &searchProducts)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "could not fetch products from cursor")
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, "could not fetch products from cursor")
		}
		return c.JSON(200, searchProducts)
	}
}
