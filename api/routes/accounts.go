package routes

import (
	"accounts.sidooh/api/middlewares"
	"accounts.sidooh/pkg"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/pkg/repositories/accounts"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)

type CreateAccountRequest struct {
	Phone      string `json:"phone" validate:"required,numeric,min=9"`
	InviteCode string `json:"invite_code,omitempty"`
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
	WithUser    string `query:"with_user" validate:"omitempty,oneof=true false"`
	WithInviter string `query:"with_inviter" validate:"omitempty,oneof=true false"`
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

// TODO: Improve error handling, statuses, messages etc...
func RegisterAccountsHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/accounts", func(context echo.Context) error {
		request := new(AccountsRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		accounts, err := accounts.GetAccounts(request.WithUser == "true", constants.DEFAULT_QUERY_LIMIT)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, accounts)
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
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		account, err := accounts.GetAccountById(uint(id), request.WithUser == "true", request.WithInviter == "true")
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/phone/:phone", func(context echo.Context) error {
		request := new(AccountByPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		// TODO: Move country to config
		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(request.Phone))
		}

		account, err := accounts.GetAccountByPhone(phone, request.WithUser == "true")
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts", func(context echo.Context) error {
		request := new(CreateAccountRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(request.Phone))
		}

		account, err := accounts.Create(entities.Account{
			Phone:      phone,
			TelcoID:    1,
			Active:     true,
			InviteCode: request.InviteCode,
		})
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/check-pin", func(context echo.Context) error {
		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(context.Param("id")))
		}

		err = accounts.CheckPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, true)
	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/set-pin", func(context echo.Context) error {
		request := new(CheckPinRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return echo.NewHTTPError(422, pkg.BadRequestError{Message: err.Error()}.Errors())
		}

		err = accounts.SetPin(uint(id), strings.TrimSpace(request.Pin))
		if err != nil {
			return echo.NewHTTPError(400, pkg.BadRequestError{Message: err.Error()}.Errors())
		}

		return utils.HandleSuccessResponse(context, true)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/search/id_or_phone", func(context echo.Context) error {
		request := new(SearchIdOrPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		if len(request.Search) >= 9 { // Most likely a phone number
			// TODO: Move country to config
			phone, err := utils.GetPhoneByCountry("KE", request.Search)
			if err == nil {
				request.Search = phone
			}
		}

		account, err := accounts.SearchByIdOrPhone(request.Search)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/search/phone", func(context echo.Context) error {
		request := new(SearchPhoneRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		accounts, err := accounts.SearchByPhone(request.Search)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, accounts)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/ancestors", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")

		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return utils.HandleErrorResponse(context, err)
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := accounts.ReadAncestors(uint(id), uint(levelLimit+1))
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/descendants", func(context echo.Context) error {
		request := new(AncestorsOrDescendantRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		levelLimit := viper.GetUint64("INVITE_LEVEL_LIMIT")
		if request.LevelLimit != "" {
			requestLevelLimit, err := strconv.ParseUint(request.LevelLimit, 10, 8)
			if err != nil {
				return utils.HandleErrorResponse(context, err)
			}
			if requestLevelLimit < levelLimit {
				levelLimit = requestLevelLimit
			}
		}

		account, err := accounts.ReadDescendants(uint(id), uint(levelLimit+1))
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.GET(constants.API_URL+"/accounts/:id/has-pin", func(context echo.Context) error {
		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		exists := accounts.HasPin(uint(id))
		if exists {
			return utils.HandleSuccessResponse(context, true)
		}

		return context.JSON(http.StatusBadRequest, utils.ErrorResponse("", false))
	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/update-profile", func(context echo.Context) error {
		request := new(UpdateProfileRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		user, err := accounts.UpdateProfile(uint(id), request.Name)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, user)
	}, authMiddleware)

	e.POST(constants.API_URL+"/accounts/:id/reset-pin", func(context echo.Context) error {
		request := new(AccountByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		id, err := strconv.ParseUint(context.Param("id"), 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		err = accounts.ResetPin(uint(id))
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, true)
	}, authMiddleware)
}
