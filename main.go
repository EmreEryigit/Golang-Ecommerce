package main

import (
	"ecomm/controller"
	"ecomm/database"
	"ecomm/middleware"
	route "ecomm/routes"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	port := "3000"

	app := controller.NewApplication(database.OpenCollection(database.Client, "Products"), database.OpenCollection(database.Client, "Users"))

	r := echo.New()
	r.Use(middleware.CurrentUser)
	userG := r.Group("/users")
	route.UserRoutes(userG)
	r.Use(middleware.Authenticate)
	r.GET("/addtocart", app.AddToCart())
	r.GET("/removeitem", app.RemoveItem())
	r.GET("/listcart", controller.GetItemFromCart())
	r.POST("/addaddress", controller.AddAddress())
	r.PUT("/edithomeaddress", controller.EditHomeAddress())
	r.PUT("/editworkaddress", controller.EditWorkAddress())
	r.GET("/deleteaddresses", controller.DeleteAddress())
	r.GET("/cartcheckout", app.BuyFromCart())
	r.GET("/instantbuy", app.InstantBuy())

	log.Fatal(r.Start(port))

}
