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
	Phone string `json:"phone" validate:"required,numeric,min=9"`
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
	AccountsRequest
	Id string `param:"id" validate:"required,numeric,min=1"`
}

type AccountByPhoneRequest struct {
	Phone string `param:"phone" validate:"required,numeric,min=9"`
	AccountsRequest
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

		accounts, err := repositories.GetAccounts(request.WithUser == "true")
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, accounts)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id", func(context echo.Context) error {
		// TODO: Review validating request using generic
		//err := middlewares.ValidateRequest(context, new(AccountByIdRequest))
		//if err != nil {
		//	return err
		//}

		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		// TODO: Refactor id/phone/etc... checks to validation framework
		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		account, err := repositories.GetAccountById(uint(id), request.WithUser == "true")
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, account)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/phone/:phone", func(context echo.Context) error {
		request := new(AccountByPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		// TODO: Move country to config
		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.PhoneValidationErrorResponse(request.Phone))
		}

		account, err := repositories.GetAccountByPhone(phone, request.WithUser == "true")
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, account)

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

		if len(request.Search) >= 9 { // Most likely a phone number
			// TODO: Move country to config
			phone, err := util.GetPhoneByCountry("KE", request.Search)
			if err == nil {
				request.Search = phone
			}
		}

		account, err := Account.SearchByIdOrPhone(request.Search)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, account)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/search/phone", func(context echo.Context) error {
		request := new(SearchPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		accounts, err := Account.SearchByPhone(request.Search)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, accounts)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/ancestors", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")

		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return util.HandleErrorResponse(context, err)
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := Account.Ancestors(uint(id), uint(levelLimit+1))
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, account)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/descendants", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")
		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return util.HandleErrorResponse(context, err)
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := Account.Descendants(uint(id), uint(levelLimit+1))
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, account)

	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/has-pin", func(context echo.Context) error {

		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		exists := repositories.HasPin(uint(id))
		if exists {
			return util.HandleSuccessResponse(context, true)
		}

		return context.JSON(http.StatusBadRequest, util.ErrorResponse("", false))

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/update-profile", func(context echo.Context) error {

		request := new(UpdateProfileRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		user, err := repositories.UpdateProfile(uint(id), request.Name)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, user)

	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/reset-pin", func(context echo.Context) error {

		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.IdValidationErrorResponse(request.Id))
		}

		err = repositories.ResetPin(uint(id))
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, true)

	}, authMiddleware)
}
