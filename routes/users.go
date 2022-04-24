package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	User "accounts.sidooh/models/user"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type SearchEmailRequest struct {
	Email string `query:"email" validate:"required,min=2,max=32"`
}

type UserByIdRequest struct {
	Id          string `param:"id" validate:"required,numeric,number,min=1"`
	WithAccount string `query:"with_account" validate:"omitempty,oneof=true false"`
}

func RegisterUsersHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {

	e.GET(constants.API_URL+"/users", func(context echo.Context) error {

		users, err := User.All()
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, users)

	}, authMiddleware)

	e.GET(constants.API_URL+"/users/:id", func(context echo.Context) error {
		request := new(UserByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		//if request.WithAccount == "true" {
		//	//TODO: To implement
		//} else {
		//	user, err := User.FindUserById(uint(id))
		//	if err != nil {
		//		return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		//	}
		//
		//	return context.JSON(http.StatusOK, user)
		//}

		user, err := User.FindUserById(uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, user)

	}, authMiddleware)

	e.GET(constants.API_URL+"/users/search", func(context echo.Context) error {
		request := new(SearchEmailRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		users, err := User.SearchByEmail(request.Email)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, users)

	}, authMiddleware)
}
