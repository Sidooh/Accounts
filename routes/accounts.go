package routes

import (
	"accounts.sidooh/db"
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Account "accounts.sidooh/models/account"
	"accounts.sidooh/repositories"
	"accounts.sidooh/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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

type AncestorsOrDescendantRequest struct {
	Id         string `param:"id" validate:"required,numeric,min=1"`
	LevelLimit string `query:"level_limit" validate:"omitempty,number,min=1,max=5"`
}

type AccountByIdRequest struct {
	Id       string `param:"id" validate:"required,numeric,min=1"`
	WithUser string `query:"with_user" validate:"omitempty,oneof=true false"`
}

type AccountByPhoneRequest struct {
	Phone    string `param:"phone" validate:"required,numeric,min=9"`
	WithUser string `query:"with_user" validate:"omitempty,oneof=true false"`
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
		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			fmt.Println(err)
			request.Id = context.Param("id")
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		if request.WithUser == "true" {
			account, err := Account.ByIdWithUser(datastore, uint(id))
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)

		} else {
			account, err := Account.ById(datastore, uint(id))
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)

		}

	}, util.CustomJWTMiddleware)

	e.GET("/api/accounts/phone/:phone", func(context echo.Context) error {
		request := new(AccountByPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		// TODO: Move country to config
		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		if request.WithUser == "true" {
			account, err := Account.ByPhoneWithUser(datastore, phone)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)
		} else {
			account, err := Account.ByPhone(datastore, phone)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)
		}

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

	e.GET("/api/accounts/:id/ancestors", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")

		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := Account.Ancestors(uint(id), uint(levelLimit+1))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)

	}, util.CustomJWTMiddleware)

	e.GET("/api/accounts/:id/descendants", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")
		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := Account.Descendants(uint(id), uint(levelLimit+1))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)

	}, util.CustomJWTMiddleware)
}
