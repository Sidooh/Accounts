package routes

import (
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterSignOutHandler(e *echo.Echo) {
	e.POST("/api/users/signout", func(context echo.Context) error {

		util.InvalidateToken(context)
		context.Set("user", nil)

		return context.JSON(http.StatusOK, context.Get("user"))
	}, util.CustomJWTMiddleware)
}
