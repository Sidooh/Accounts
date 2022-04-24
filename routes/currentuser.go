package routes

import (
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	//TODO(Check an update to jwt of echo using V4, currently on v3
	"github.com/golang-jwt/jwt"
	"net/http"
)

func RegisterCurrentUserHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/users/currentuser", func(context echo.Context) error {

		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(*util.MyCustomClaims)

		return context.JSON(http.StatusOK, claims)

	}, authMiddleware)
}
