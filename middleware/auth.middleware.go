package middleware

import (
	"ecomm/controller"
	"ecomm/helper"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("current-user").(*helper.SignedDetails)
		if currentUser == nil {
			return c.JSON(http.StatusUnauthorized, "Restricted route")
		}
		return next(c)
	}
}

func CurrentUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := controller.Store.Get(c.Request(), "auth-session")
		strJwt := fmt.Sprintf("%v", session.Values["auth"])
		claims, _ := helper.ValidateToken(strJwt)
		c.Set("current-user", claims)
		return next(c)
	}
}
