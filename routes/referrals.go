package routes

import (
	"accounts.sidooh/db"
	"accounts.sidooh/errors"
	"accounts.sidooh/middlewares"
	Referral "accounts.sidooh/models/referral"
	"accounts.sidooh/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateReferralRequest struct {
	AccountId     uint   `json:"account_id" form:"account_id" validate:"required,numeric"`
	ReferralPhone string `json:"referral_phone" form:"referral_phone" validate:"required,numeric"`
}

func RegisterReferralsHandler(e *echo.Echo) {
	e.GET("/api/referrals", func(context echo.Context) error {

		datastore := db.NewConnection()
		referrals, err := Referral.All(datastore)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referrals)

	}, util.CustomJWTMiddleware)

	e.GET("/api/referrals/:phone", func(context echo.Context) error {

		phone := context.Param("phone")

		phone, err := util.GetPhoneByCountry("KE", phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		referral, err := Referral.ByPhone(phone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referral)

	}, util.CustomJWTMiddleware)

	e.POST("/api/referrals", func(context echo.Context) error {

		request := new(CreateReferralRequest)
		if err := middlewares.BindAndValidateRequest(context, request); err != nil {
			return err
		}

		phone, err := util.GetPhoneByCountry("KE", request.ReferralPhone)
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		//accountId, err := strconv.ParseUint(request.AccountId, 10, 32)
		//if err != nil {
		//	return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		//}

		datastore := db.NewConnection()
		referral, err := Referral.Create(datastore, Referral.Model{
			AccountID:    request.AccountId,
			RefereePhone: phone,
		})
		if err != nil {
			return echo.NewHTTPError(400, errors.BadRequestError{Message: err.Error()}.Errors())
		}

		return context.JSON(http.StatusOK, referral)
	})
}
