package routes

import (
	"accounts.sidooh/api/middlewares"
	"accounts.sidooh/pkg/entities"
	invitesRepo "accounts.sidooh/pkg/repositories/invites"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateInviteRequest struct {
	InviterId string `json:"inviter_id" form:"inviter_id" validate:"required,numeric"`
	Phone     string `json:"phone" form:"phone" validate:"required,numeric"`
	Type      string `query:"type" validate:"omitempty"`
}

type InvitesRequest struct {
	With string `query:"with" validate:"omitempty,containsany=account inviter"`
	Days string `query:"days" validate:"omitempty,numeric"`
}

type InviteByIdRequest struct {
	InvitesRequest
	Id string `param:"id" validate:"required,numeric,min=1"`
}

func RegisterInvitesHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/invites", func(context echo.Context) error {
		request := new(InvitesRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		days, _ := strconv.ParseUint(request.Days, 10, 32)

		invites, err := invitesRepo.GetInvites(request.With, int(days), constants.DEFAULT_QUERY_LIMIT)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, invites)
	}, authMiddleware)

	e.GET(constants.API_URL+"/invites/:id", func(context echo.Context) error {
		// TODO: Review validating request using generic
		request := new(InviteByIdRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			request.Id = context.Param("id")
		}

		// TODO: Refactor id/phone/etc... checks to validation framework
		id, err := strconv.ParseUint(request.Id, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.IdValidationErrorResponse(request.Id))
		}

		account, err := invitesRepo.GetInviteById(uint(id), request.With)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, account)
	}, authMiddleware)

	e.GET(constants.API_URL+"/invites/phone/:phone", func(context echo.Context) error {
		phone := context.Param("phone")

		phone, err := utils.GetPhoneByCountry("KE", phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(phone))
		}

		invite, err := invitesRepo.ReadByPhone(phone)
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, invite)
	}, authMiddleware)

	e.POST(constants.API_URL+"/invites", func(context echo.Context) error {
		request := new(CreateInviteRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		//TODO: Move Country to env and fetch from it
		phone, err := utils.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.PhoneValidationErrorResponse(request.Phone))
		}

		InviterId, err := strconv.ParseUint(request.InviterId, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, utils.InviterIdValidationErrorResponse(request.InviterId))
		}

		invite, err := invitesRepo.Create(entities.Invite{
			InviterID: uint(InviterId),
			Phone:     phone,
			Type:      request.Type,
		})
		if err != nil {
			return utils.HandleErrorResponse(context, err)
		}

		return utils.HandleSuccessResponse(context, invite)
	}, authMiddleware)
}
