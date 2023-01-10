package routes

import (
	"accounts.sidooh/middlewares"
	Invite "accounts.sidooh/models/invite"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CreateInviteRequest struct {
	InviterId string `json:"inviter_id" form:"inviter_id" validate:"required,numeric"`
	Phone     string `json:"phone" form:"phone" validate:"required,numeric"`
}

func RegisterInvitesHandler(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	e.GET(constants.API_URL+"/invites", func(context echo.Context) error {

		invites, err := Invite.All(constants.DEFAULT_QUERY_LIMIT)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, invites)

	}, authMiddleware)

	e.GET(constants.API_URL+"/invites/phone/:phone", func(context echo.Context) error {

		phone := context.Param("phone")

		phone, err := util.GetPhoneByCountry("KE", phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.PhoneValidationErrorResponse(phone))
		}

		invite, err := Invite.ByPhone(phone)
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, invite)
	}, authMiddleware)

	e.POST(constants.API_URL+"/invites", func(context echo.Context) error {
		request := new(CreateInviteRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		//TODO: Move Country to env and fetch from it
		phone, err := util.GetPhoneByCountry("KE", request.Phone)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.PhoneValidationErrorResponse(request.Phone))
		}

		InviterId, err := strconv.ParseUint(request.InviterId, 10, 32)
		if err != nil {
			return context.JSON(http.StatusUnprocessableEntity, util.InviterIdValidationErrorResponse(request.InviterId))
		}

		invite, err := Invite.Create(Invite.Model{
			InviterID: uint(InviterId),
			Phone:     phone,
		})
		if err != nil {
			return util.HandleErrorResponse(context, err)
		}

		return util.HandleSuccessResponse(context, invite)
	}, authMiddleware)
}
