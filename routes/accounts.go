package routes

import (
	"accounts.sidooh/db"
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Account "accounts.sidooh/models/account"
	"accounts.sidooh/repositories"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type CreateAccountRequest struct {
	Phone string `json:"phone" validate:"required,numeric"`
}

type CheckPinRequest struct {
	Pin string `json:"pin" validate:"required,numeric,min=4,max=4"`
}

type SearchPhoneRequest struct {
	Phone string `query:"phone" validate:"required,numeric,min=3,max=12"`
}

func RegisterAccountsHandler(e *echo.Echo) {
	// TODO: Refactor these; move to repo. Repo should determine and setup datastore independently
	datastore := db.NewConnection()
	repositories.Construct(datastore)

	e.GET("/api/accounts", func(context echo.Context) error {

		accounts, err := Account.All(datastore)
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

		account, err := Account.ById(datastore, uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)

	}, util.CustomJWTMiddleware)

	e.GET("/api/accounts/phone/:phone", func(context echo.Context) error {

		phone, err := util.GetPhoneByCountry("KE", context.Param("phone"))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		account, err := Account.ByPhone(datastore, phone)
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

	e.POST("/api/accounts/:id/check-pin", func(context echo.Context) error {

		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.CheckPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})

	}, util.CustomJWTMiddleware)

	e.POST("/api/accounts/:id/set-pin", func(context echo.Context) error {
		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.SetPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	}, util.CustomJWTMiddleware)

	e.GET("/api/accounts/search", func(context echo.Context) error {
		request := new(SearchPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		accounts, err := Account.SearchByPhone(datastore, request.Phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, accounts)

	}, util.CustomJWTMiddleware)
}
