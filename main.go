package main

import (
	"ecomm/controller"
	"ecomm/database"
	"log"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/route"
)

func main() {
	port := "3000"

	app := controller.NewApplication(database.OpenCollection(database.Client, "Products"), database.OpenCollection(database.Client, "Users"))

	r := echo.New()
	userG := r.Group("/users")
	route.UserRoutes(userG)

	r.GET("/addtocart")
	r.GET("/removeitem")
	r.GET("/cartcheckout")
	r.GET("/instantbuy")

	log.Fatal(r.Start(port))

}
