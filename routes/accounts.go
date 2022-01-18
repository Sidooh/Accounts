package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Account "accounts.sidooh/models/account"
	"accounts.sidooh/models/repositories"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateAccountRequest struct {
	Phone string `json:"phone" form:"phone" validate:"required,numeric"`
}

func RegisterAccountsHandler(e *echo.Echo) {
	e.GET("/api/accounts", func(context echo.Context) error {

		accounts, err := Account.All()
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, accounts)

	}, util.CustomJWTMiddleware)

	e.GET("/api/accounts/:id", func(context echo.Context) error {

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		account, err := Account.ById(uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)

	}, util.CustomJWTMiddleware)

	e.POST("/api/accounts", func(context echo.Context) error {

		request := new(CreateAccountRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		account, err := repositories.Create(Account.Model{
			Phone:   phone,
			TelcoID: 1,
			Active:  true,
		})
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)
	})
}
