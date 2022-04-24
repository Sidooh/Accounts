package routes

import (
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterSignOutHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.POST(constants.API_URL+"/users/signout", func(context echo.Context) error {

		util.InvalidateToken(context)
		context.Set("user", nil)

		return context.JSON(http.StatusOK, context.Get("user"))
	}, authMiddleware)
}
