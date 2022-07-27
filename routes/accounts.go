package routes

import (
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Account "accounts.sidooh/models/account"
	"accounts.sidooh/repositories"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
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
	Search string `query:"search" validate:"required,numeric,min=3,max=12"`
}

type SearchIdOrPhoneRequest struct {
	Search string `query:"search" validate:"required,numeric,min=1,max=12"`
}

type AncestorsOrDescendantRequest struct {
	Id         string `param:"id" validate:"required,numeric,min=1"`
	LevelLimit string `query:"level_limit" validate:"omitempty,number,min=1,max=5"`
}

type AccountsRequest struct {
	WithUser string `query:"with_user" validate:"omitempty,oneof=true false"`
}

type AccountByIdRequest struct {
	Id       string `param:"id" validate:"required,numeric,min=1"`
	WithUser string `query:"with_user" validate:"omitempty,oneof=true false"`
}

type AccountByPhoneRequest struct {
	Phone    string `param:"phone" validate:"required,numeric,min=9"`
	WithUser string `query:"with_user" validate:"omitempty,oneof=true false"`
}

type UpdateProfileRequest struct {
	Id   string `param:"id" validate:"required,numeric,min=1"`
	Name string `json:"name" validate:"required,min=3"`
}

//TODO: Improve error handling, statuses, messages etc...
func RegisterAccountsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/accounts", func(context echo.Context) error {
		request := new(AccountsRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		if request.WithUser == "true" {
			accounts, err := Account.AllWithUser()
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, accounts)

		} else {
			accounts, err := Account.All()
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, accounts)

		}

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id", func(context echo.Context) error {
		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
			//return echo.NewHTTPError(422, errors.ValidationError{
			//	Value:   request.Id,
			//	Message: "valid id is required",
			//	Param:   "id",
			//})
		}

		if request.WithUser == "true" {
			account, err := Account.ByIdWithUser(uint(id))
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)

		} else {
			account, err := Account.ById(uint(id))
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)

		}

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/phone/:phone", func(context echo.Context) error {
		request := new(AccountByPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		// TODO: Move country to config
		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		if request.WithUser == "true" {
			account, err := Account.ByPhoneWithUser(phone)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)
		} else {
			account, err := Account.ByPhone(phone)
			if err != nil {
				return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
			}

			return context.JSON(http.StatusOK, account)
		}

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts", func(context echo.Context) error {
		request := new(CreateAccountRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
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

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/check-pin", func(context echo.Context) error {

		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.CheckPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/set-pin", func(context echo.Context) error {
		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.SetPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/search/id_or_phone", func(context echo.Context) error {
		request := new(SearchIdOrPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		account, err := Account.SearchByIdOrPhone(request.Search)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, account)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/search/phone", func(context echo.Context) error {
		request := new(SearchPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		accounts, err := Account.SearchByPhone(request.Search)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, accounts)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/ancestors", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")

		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
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

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/descendants", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")
		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
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

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/has-pin", func(context echo.Context) error {

		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		err = repositories.HasPin(uint(id))
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, map[string]bool{
			"message": true,
		})

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/update-profile", func(context echo.Context) error {

		request := new(UpdateProfileRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		user, err := repositories.UpdateProfile(uint(id), request.Name)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, user)

	}, authMiddleware)

}
