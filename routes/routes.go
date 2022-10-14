package routes

import (
	"ecomm/controller"

	"github.com/labstack/echo/v4"
)

func UserRoutes(g *echo.Group) {
	g.POST("/users/signup", controller.Signup())
	g.POST("/users/login", controller.Login())
	g.POST("/admin/addproduct", controller.ProductViewerAdmin())
	g.GET("/users/productview", controller.SearchProduct())
	g.GET("/users/search", controller.SearchProductByQuery())
}
