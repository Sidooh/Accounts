package routes

import (
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterSignOutHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.POST(constants.API_URL+"/users/signout", func(context echo.Context) error {
		utils.InvalidateToken(context)
		context.Set("user", nil)

		return context.JSON(http.StatusOK, context.Get("user"))
	}, authMiddleware)
}
