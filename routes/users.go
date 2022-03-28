package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SearchEmailRequest struct {
	Email string `query:"email" validate:"required,min=2,max=32"`
}

func RegisterUsersHandler(e *echo.Echo) {

	e.GET("/api/users", func(context echo.Context) error {

		users, err := User.All()
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, users)

	}, util.CustomJWTMiddleware)

	e.GET("/api/users/search", func(context echo.Context) error {
		request := new(SearchEmailRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		users, err := User.SearchByEmail(request.Email)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, users)

	}, util.CustomJWTMiddleware)
}
